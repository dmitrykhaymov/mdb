package api

import (
	"database/sql"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/queries/qm"
	"gopkg.in/nullbio/null.v6"

	"encoding/json"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
	"strings"
)

// Do all stuff for processing metadata coming from Content Identification Tool.
// 	1. Update properties for original and proxy (film_date, capture_date)
//	2. Update language of original
// 	3. Create content_unit (content_type, dates)
//	4. Add files to new unit
// 	5. Add ancestor files to unit
// 	6. Associate unit with sources, tags, and persons
// 	7. Get or create collection
// 	8. Update collection (content_type, dates, number) if full lesson or new lesson
// 	9. Associate collection and unit
// 	10. Associate unit and derived units
// 	11. Set default permissions ?!
func ProcessCITMetadata(exec boil.Executor, metadata CITMetadata, original, proxy *models.File) error {
	log.Info("CITMetadata: %v", metadata)

	// Update properties for original and proxy (film_date, capture_date)
	filmDate := metadata.CaptureDate
	if metadata.WeekDate != nil {
		filmDate = *metadata.WeekDate
	}
	props := map[string]interface{}{
		"capture_date": metadata.CaptureDate,
		"film_date":    filmDate,
	}
	err := UpdateFileProperties(exec, original, props)
	if err != nil {
		return err
	}
	err = UpdateFileProperties(exec, proxy, props)
	if err != nil {
		return err
	}

	// Update language of original.
	// TODO: What about proxy !?
	if metadata.HasTranslation {
		original.Language = null.StringFrom(LANG_MULTI)
	} else {
		l := StdLang(metadata.Language)
		if l == LANG_UNKNOWN {
			log.Warnf("Unknown language in metadata %s", metadata.Language)
		}
		original.Language = null.StringFrom(l)
	}
	err = original.Update(exec, "language")
	if err != nil {
		return errors.Wrap(err, "Save original to DB")
	}

	// Create content_unit (content_type, dates)
	cu, err := CreateContentUnit(exec, metadata.ContentType, props)
	if err != nil {
		return errors.Wrap(err, "Create content unit")
	}

	// Add files to new unit
	err = cu.AddFiles(exec, false, original, proxy)
	if err != nil {
		return errors.Wrap(err, "Add files to unit")
	}

	// Add ancestor files to unit (not for derived units)
	if !metadata.ArtifactType.Valid ||
		metadata.ArtifactType.String == "main" {
		ancestors, err := FindFileAncestors(exec, original.ID)
		if err != nil {
			return errors.Wrap(err, "Find original's ancestors")
		}
		err = cu.AddFiles(exec, false, ancestors...)
		if err != nil {
			return errors.Wrap(err, "Add ancestors to unit")
		}
	}

	// Associate unit with sources, tags, and persons
	if len(metadata.Sources) > 0 {
		sources, err := models.Sources(exec,
			qm.WhereIn("uid in ?", utils.ConvertArgsString(metadata.Sources)...)).
			All()
		if err != nil {
			return errors.Wrap(err, "Lookup sources in DB")
		}
		if len(sources) != len(metadata.Sources) {
			missing := make([]string, 0)
			for _, x := range metadata.Sources {
				found := false
				for _, y := range sources {
					if x == y.UID {
						found = true
						break
					}
				}
				if !found {
					missing = append(missing, x)
				}
			}
			log.Warnf("Unknown sources: %s", missing)
		}
		err = cu.AddSources(exec, false, sources...)
		if err != nil {
			return errors.Wrap(err, "Associate sources")
		}
	}

	if len(metadata.Tags) > 0 {
		tags, err := models.Tags(exec,
			qm.WhereIn("uid in ?", utils.ConvertArgsString(metadata.Tags)...)).
			All()
		if err != nil {
			return errors.Wrap(err, "Lookup tags  in DB")
		}
		if len(tags) != len(metadata.Tags) {
			missing := make([]string, 0)
			for _, x := range metadata.Tags {
				found := false
				for _, y := range tags {
					if x == y.UID {
						found = true
						break
					}
				}
				if !found {
					missing = append(missing, x)
				}
			}
			log.Warnf("Unknown sources: %s", missing)
		}
		err = cu.AddTags(exec, false, tags...)
		if err != nil {
			return errors.Wrap(err, "Associate tags")
		}
	}

	// Handle persons ...
	if strings.ToLower(metadata.Lecturer) == P_RAV {
		cup := &models.ContentUnitsPerson{
			PersonID: PERSONS_REGISTRY.ByPattern[P_RAV].ID,
			RoleID: CONTENT_ROLE_TYPE_REGISTRY.ByName[CR_LECTURER].ID,
		}
		err = cu.AddContentUnitsPersons(exec, true, cup)
		if err != nil {
			return errors.Wrap(err, "Associate persons")
		}
	}


	// Get or create collection
	var c *models.Collection
	if metadata.CollectionUID.Valid {
		c, err = models.Collections(exec, qm.Where("uid = ?", metadata.CollectionUID.String)).One()
		if err != nil {
			if err == sql.ErrNoRows {
				return errors.Wrapf(err, "No such collection, uid=%s",
					metadata.CollectionUID.String)
			} else {
				return errors.Wrap(err, "Lookup collection in DB")
			}
		}
	} else if metadata.ContentType == CT_LESSON_PART ||
		metadata.ContentType == CT_FULL_LESSON {

		// Reconcile or create new
		// Reconciliation is done by looking up the operation chain of original to capture_stop.
		// There we have a property of saying the capture_id of the full lesson capture.
		captureStop, err := FindUpChainOperation(exec, original.ID,
			OPERATION_TYPE_REGISTRY.ByName[OP_CAPTURE_STOP].ID)
		if err != nil {
			if ex, ok := err.(UpChainOperationNotFound); ok {
				log.Warnf(ex.Error())
			} else {
				return err
			}
		} else if captureStop.Properties.Valid {
			var oProps map[string]interface{}
			err = json.Unmarshal(captureStop.Properties.JSON, &oProps)
			if err != nil {
				return errors.Wrap(err, "json Unmarshal")
			}
			captureID, ok := oProps["collection_uid"]
			if ok {
				var ct string
				if metadata.WeekDate == nil {
					ct = CT_DAILY_LESSON
				} else {
					ct = CT_SATURDAY_LESSON
				}

				// Keep this property on the collection for other parts to find it
				props["capture_id"] = captureID
				if metadata.Number.Valid {
					props["number"] = metadata.Number.Int
				}

				c, err = FindCollectionByCaptureID(exec, captureID)
				if err != nil {
					if _, ok := err.(CollectionNotFound); !ok {
						return err
					}

					// Create new collection
					c, err = CreateCollection(exec, ct, props)
					if err != nil {
						return err
					}
				} else if metadata.ContentType == CT_FULL_LESSON {
					// Update collection properties to those of full lesson
					if c.TypeID != CONTENT_TYPE_REGISTRY.ByName[ct].ID {
						c.TypeID = CONTENT_TYPE_REGISTRY.ByName[ct].ID
						err = c.Update(exec, "type_id")
						if err != nil {
							return errors.Wrap(err, "Update collection type in DB")
						}
					}

					err = UpdateCollectionProperties(exec, c, props)
					if err != nil {
						return err
					}
				}
			} else {
				log.Warnf("No collection_uid in capture_stop [%d] properties", captureStop.ID)
			}
		} else {
			log.Warnf("Invalid properties in capture_stop [%d]", captureStop.ID)
		}
	}

	// Associate collection and unit
	if c != nil &&
		(!metadata.ArtifactType.Valid || metadata.ArtifactType.String == "main") {
		ccu := &models.CollectionsContentUnit{
			CollectionID:  c.ID,
			ContentUnitID: cu.ID,
		}
		switch metadata.ContentType {
		case CT_FULL_LESSON:
			if c.TypeID == CONTENT_TYPE_REGISTRY.ByName[CT_DAILY_LESSON].ID ||
				c.TypeID == CONTENT_TYPE_REGISTRY.ByName[CT_SATURDAY_LESSON].ID {
				ccu.Name = "full"
			} else if metadata.Number.Valid {
				ccu.Name = strconv.Itoa(metadata.Number.Int)
			}
			break
		case CT_LESSON_PART:
			if metadata.Part.Valid {
				ccu.Name = strconv.Itoa(metadata.Part.Int)
			}
			break
		case CT_VIDEO_PROGRAM_CHAPTER:
			if metadata.Episode.Valid {
				ccu.Name = metadata.Episode.String
			}
			break
		default:
			if metadata.Number.Valid {
				ccu.Name = strconv.Itoa(metadata.Number.Int)
			}
			if metadata.PartType.Valid && metadata.PartType.Int > 2 {
				idx := metadata.PartType.Int - 3
				if idx < len(MISC_EVENT_PART_TYPES) {
					ccu.Name = MISC_EVENT_PART_TYPES[idx] + ccu.Name
				} else {
					log.Warn("Unknown event part type: %d", metadata.PartType.Int)
				}
			}
			break
		}

		err = c.AddCollectionsContentUnits(exec, true, ccu)
		if err != nil {
			return errors.Wrap(err, "Save collection and content unit association in DB")
		}
	}

	// associate unit and derived units

	// set default permissions ?!

	return nil
}
