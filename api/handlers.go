package api

import (
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/queries"
	"github.com/vattle/sqlboiler/queries/qm"
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/nullbio/null.v6"

	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
)

// Start capture of AV file, i.e. morning lesson, tv program, etc...
func CaptureStartHandler(c *gin.Context) {
	log.Info(OP_CAPTURE_START)
	var i CaptureStartRequest
	if c.BindJSON(&i) == nil {
		handleOperation(c, i, handleCaptureStart)
	}
}

// Stop capture of AV file, ending a matching capture_start event.
// This is the first time a physical file is created in the studio.
func CaptureStopHandler(c *gin.Context) {
	log.Info(OP_CAPTURE_STOP)
	var i CaptureStopRequest
	if c.BindJSON(&i) == nil {
		handleOperation(c, i, handleCaptureStop)
	}
}

// Demux manifest file to original and low resolution proxy
func DemuxHandler(c *gin.Context) {
	log.Info(OP_DEMUX)
	var i DemuxRequest
	if c.BindJSON(&i) == nil {
		handleOperation(c, i, handleDemux)
	}
}

// Trim demuxed files at certain points
func TrimHandler(c *gin.Context) {
	log.Info(OP_TRIM)
	var i TrimRequest
	if c.BindJSON(&i) == nil {
		handleOperation(c, i, handleTrim)
	}
}

// Final files sent from studio
func SendHandler(c *gin.Context) {
	log.Info(OP_SEND)
	var i SendRequest
	if c.BindJSON(&i) == nil {
		handleOperation(c, i, handleSend)
	}
}

// Files converted to low resolution web formats, language splitting, etc...
func ConvertHandler(c *gin.Context) {
	log.Info(OP_CONVERT)
	var i ConvertRequest
	if c.BindJSON(&i) == nil {
		handleOperation(c, i, handleConvert)
	}
}

// File uploaded to a public accessible URL
func UploadHandler(c *gin.Context) {
	log.Info(OP_UPLOAD)
	var i UploadRequest
	if c.BindJSON(&i) == nil {
		handleOperation(c, i, handleUpload)
	}
}

// Handler logic

func handleCaptureStart(exec boil.Executor, input interface{}) (*models.Operation, error) {
	r := input.(CaptureStartRequest)

	log.Info("Creating operation")
	props := map[string]interface{}{
		"capture_source": r.CaptureSource,
		"collection_uid": r.CollectionUID,
	}
	operation, err := CreateOperation(exec, OP_CAPTURE_START, r.Operation, props)
	if err != nil {
		return nil, err
	}

	log.Info("Creating file and associating to operation")
	file := models.File{
		UID:  utils.GenerateUID(8),
		Name: r.FileName,
	}
	return operation, operation.AddFiles(exec, true, &file)
}

func handleCaptureStop(exec boil.Executor, input interface{}) (*models.Operation, error) {
	r := input.(CaptureStopRequest)

	log.Info("Creating operation")
	props := map[string]interface{}{
		"capture_source": r.CaptureSource,
		"collection_uid": r.CollectionUID, // $LID = backup capture id when lesson, capture_id when program (part=false)
		"part":           r.Part,
	}
	operation, err := CreateOperation(exec, OP_CAPTURE_STOP, r.Operation, props)
	if err != nil {
		return nil, err
	}

	log.Info("Looking up parent file, workflow_id=", r.Operation.WorkflowID)
	var parent *models.File
	var parentID int64
	err = queries.Raw(exec,
		`SELECT file_id FROM files_operations
		 INNER JOIN operations ON operation_id = id
		 WHERE type_id=$1 AND properties -> 'workflow_id' ? $2`,
		OPERATION_TYPE_REGISTRY.ByName[OP_CAPTURE_START].ID,
		r.Operation.WorkflowID).
		QueryRow().
		Scan(&parentID)
	if err == nil {
		parent = &models.File{ID: parentID}
	} else {
		if err == sql.ErrNoRows {
			log.Warnf("capture_start operation not found for workflow_id [%s]. Skipping.",
				r.Operation.WorkflowID)
		} else {
			return nil, err
		}
	}

	log.Info("Creating file")
	file, err := CreateFile(exec, parent, r.File, nil)
	if err != nil {
		return nil, err
	}

	log.Info("Associating file to operation")
	return operation, operation.AddFiles(exec, false, file)
}

func handleDemux(exec boil.Executor, input interface{}) (*models.Operation, error) {
	r := input.(DemuxRequest)

	parent, _, err := FindFileBySHA1(exec, r.Sha1)
	if err != nil {
		return nil, err
	}

	log.Info("Creating operation")
	props := map[string]interface{}{
		"capture_source": r.CaptureSource,
	}
	operation, err := CreateOperation(exec, OP_DEMUX, r.Operation, props)
	if err != nil {
		return nil, err
	}

	log.Info("Creating original")
	props = map[string]interface{}{
		"duration": r.Original.Duration,
	}
	original, err := CreateFile(exec, parent, r.Original.File, props)
	if err != nil {
		return nil, err
	}

	log.Info("Creating proxy")
	props = map[string]interface{}{
		"duration": r.Proxy.Duration,
	}
	proxy, err := CreateFile(exec, parent, r.Proxy.File, props)
	if err != nil {
		return nil, err
	}

	log.Info("Associating files to operation")
	return operation, operation.AddFiles(exec, false, parent, original, proxy)
}

func handleTrim(exec boil.Executor, input interface{}) (*models.Operation, error) {
	r := input.(TrimRequest)

	// Fetch parent files
	original, _, err := FindFileBySHA1(exec, r.OriginalSha1)
	if err != nil {
		return nil, err
	}
	proxy, _, err := FindFileBySHA1(exec, r.ProxySha1)
	if err != nil {
		return nil, err
	}

	// TODO: in case of re-trim with the exact same parameters we already have the files in DB.
	// No need to return an error, a warning in the log is enough.

	log.Info("Creating operation")
	props := map[string]interface{}{
		"capture_source": r.CaptureSource,
		"in":             r.In,
		"out":            r.Out,
	}
	operation, err := CreateOperation(exec, OP_TRIM, r.Operation, props)
	if err != nil {
		return nil, err
	}

	log.Info("Creating trimmed original")
	props = map[string]interface{}{
		"duration": r.Original.Duration,
	}
	originalTrim, err := CreateFile(exec, original, r.Original.File, props)
	if err != nil {
		return nil, err
	}

	log.Info("Creating trimmed proxy")
	props = map[string]interface{}{
		"duration": r.Proxy.Duration,
	}
	proxyTrim, err := CreateFile(exec, proxy, r.Proxy.File, props)
	if err != nil {
		return nil, err
	}

	log.Info("Associating files to operation")
	return operation, operation.AddFiles(exec, false, original, originalTrim, proxy, proxyTrim)
}

func handleSend(exec boil.Executor, input interface{}) (*models.Operation, error) {
	r := input.(SendRequest)

	// Original
	original, _, err := FindFileBySHA1(exec, r.Original.Sha1)
	if err != nil {
		return nil, err
	}
	if original.Name == r.Original.FileName {
		log.Info("Original's name hasn't change")
	} else {
		log.Info("Renaming original")
		original.Name = r.Original.FileName
		err = original.Update(exec, "name")
		if err != nil {
			return nil, err
		}
	}

	// Proxy
	proxy, _, err := FindFileBySHA1(exec, r.Proxy.Sha1)
	if err != nil {
		return nil, err
	}
	if proxy.Name == r.Proxy.FileName {
		log.Info("Proxy's name hasn't change")
	} else {
		log.Info("Renaming proxy")
		proxy.Name = r.Proxy.FileName
		err = proxy.Update(exec, "name")
		if err != nil {
			return nil, err
		}
	}

	log.Info("Creating operation")
	var props map[string]interface{}
	if r.Metadata != nil {
		b, err := json.Marshal(r.Metadata)
		if err == nil {
			return nil, errors.Wrap(err, "json Marshal CITMetadata")
		}
		if err = json.Unmarshal(b, &props); err != nil {
			return nil, errors.Wrap(err, "json Unmarshal CITMetadata")
		}
	}
	operation, err := CreateOperation(exec, OP_SEND, r.Operation, props)
	if err != nil {
		return nil, err
	}

	log.Info("Associating files to operation")
	return operation, operation.AddFiles(exec, false, original, proxy)
}

func handleConvert(exec boil.Executor, input interface{}) (*models.Operation, error) {
	r := input.(ConvertRequest)

	in, _, err := FindFileBySHA1(exec, r.Sha1)
	if err != nil {
		return nil, err
	}

	log.Info("Creating operation")
	operation, err := CreateOperation(exec, OP_CONVERT, r.Operation, nil)
	if err != nil {
		return nil, err
	}

	log.Info("Creating output files")
	files := make([]*models.File, len(r.Output)+1)
	files[0] = in
	props := make(map[string]interface{})
	for i, x := range r.Output {
		props["duration"] = x.Duration
		f, err := CreateFile(exec, in, x.File, props)
		if err == nil {
			files[i+1] = f
		} else {
			return nil, err
		}
	}

	log.Info("Associating files to operation")
	return operation, operation.AddFiles(exec, false, files...)
}

func handleUpload(exec boil.Executor, input interface{}) (*models.Operation, error) {
	r := input.(UploadRequest)

	log.Info("Creating operation")
	operation, err := CreateOperation(exec, OP_UPLOAD, r.Operation, nil)
	if err != nil {
		return nil, err
	}

	file, _, err := FindFileBySHA1(exec, r.Sha1)
	if err != nil {
		if _, ok := err.(FileNotFound); ok {
			log.Info("File not found, creating new.")
			file, err = CreateFile(exec, nil, r.File, nil)
		} else {
			return nil, err
		}
	}

	log.Info("Updating file's properties")
	var fileProps = make(map[string]interface{})
	if file.Properties.Valid {
		file.Properties.Unmarshal(&fileProps)
	}
	fileProps["url"] = r.Url
	fileProps["duration"] = r.Duration
	fpa, _ := json.Marshal(fileProps)
	file.Properties = null.JSONFrom(fpa)

	log.Info("Saving changes to DB")
	err = file.Update(exec, "properties")
	if err != nil {
		return nil, err
	}

	log.Info("Associating file to operation")
	return operation, operation.AddFiles(exec, false, file)
}

// Helpers

// Generic operation handler.
// 	* Manage DB transactions
// 	* Call operation logic handler
// 	* Render JSON response
func handleOperation(c *gin.Context, input interface{},
	opHandler func(boil.Executor, interface{}) (*models.Operation, error)) {

	tx, err := boil.Begin()
	utils.Must(err)

	_, err = opHandler(tx, input)
	if err == nil {
		utils.Must(tx.Commit())
	} else {
		utils.Must(tx.Rollback())
		err = errors.Wrapf(err, "Handle operation")
	}

	if err == nil {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	} else {
		switch err.(type) {
		case FileNotFound:
			NewBadRequestError(err).Abort(c)
		default:
			NewInternalError(err).Abort(c)
		}
	}
}

func CreateOperation(exec boil.Executor, name string, o Operation, properties map[string]interface{}) (*models.Operation, error) {
	operation := models.Operation{
		TypeID:  OPERATION_TYPE_REGISTRY.ByName[name].ID,
		UID:     utils.GenerateUID(8),
		Station: null.StringFrom(o.Station),
	}

	// Lookup user, skip if doesn't exist
	user, err := models.Users(exec, qm.Where("email=?", o.User)).One()
	if err == nil {
		operation.UserID = null.Int64From(user.ID)
	} else {
		if err == sql.ErrNoRows {
			log.Debugf("Unknown User [%s]. Skipping.", o.User)
		} else {
			return nil, err
		}
	}

	// Handle properties
	if o.WorkflowID != "" {
		if properties == nil {
			properties = make(map[string]interface{})
		}
		properties["workflow_id"] = o.WorkflowID
	}
	if properties != nil {
		props, err := json.Marshal(properties)
		if err != nil {
			return nil, err
		}
		operation.Properties = null.JSONFrom(props)
	}

	return &operation, operation.Insert(exec)
}

func CreateFile(exec boil.Executor, parent *models.File, f File, properties map[string]interface{}) (*models.File, error) {
	sha1, err := hex.DecodeString(f.Sha1)
	if err != nil {
		return nil, err
	}

	// validate language + convert from code3 if necessary
	var mdbLang = ""
	switch len(f.Language) {
	case 2:
		if !KNOWN_LANGS.MatchString(strings.ToLower(f.Language)) {
			return nil, errors.Errorf("Unknown language %s", f.Language)
		}
		mdbLang = strings.ToLower(f.Language)
		break
	case 3:
		var ok bool
		mdbLang, ok = LANG_MAP[strings.ToUpper(f.Language)]
		if !ok {
			return nil, errors.Errorf("Unknown language %s", f.Language)
		}
		break
	case 0:
		break
	default:
		return nil, errors.Errorf("Unknown language %s", f.Language)
	}

	file := models.File{
		UID:           utils.GenerateUID(8),
		Name:          f.FileName,
		Sha1:          null.BytesFrom(sha1),
		Size:          f.Size,
		FileCreatedAt: null.TimeFrom(f.CreatedAt.Time),
		Type:          f.Type,
		SubType:       f.SubType,
		Language:      null.NewString(mdbLang, mdbLang != ""),
	}

	if f.MimeType != "" {
		file.MimeType = null.StringFrom(f.MimeType)

		// Try to complement missing type and subtype
		if file.Type == "" && file.SubType == "" {
			if mt, ok := MEDIA_TYPE_REGISTRY.ByMime[strings.ToLower(f.MimeType)]; ok {
				file.Type = mt.Type
				file.SubType = mt.SubType
			}
		}
	}

	if parent != nil {
		file.ParentID = null.Int64From(parent.ID)
	}

	// Handle properties
	if properties != nil {
		props, err := json.Marshal(properties)
		if err != nil {
			return nil, err
		}
		file.Properties = null.JSONFrom(props)
	}

	return &file, file.Insert(exec)
}

func FindFileBySHA1(exec boil.Executor, sha1 string) (*models.File, []byte, error) {
	log.Debugf("Looking up file, sha1=%s", sha1)
	s, err := hex.DecodeString(sha1)
	if err != nil {
		return nil, nil, err
	}

	f, err := models.Files(exec, qm.Where("sha1=?", s)).One()
	if err == nil {
		return f, s, nil
	} else {
		if err == sql.ErrNoRows {
			return nil, s, FileNotFound{Sha1: sha1}
		} else {
			return nil, s, err
		}
	}
}
