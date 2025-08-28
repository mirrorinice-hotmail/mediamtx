package core

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	goversion "github.com/hashicorp/go-version"
	"github.com/liip/sheriff"
)

const STREAM_LIST_FILENAME = "stream_list.json"
const STREAM_LIST_REMOTE_DB_FILENAME = "stream_list_from_db.json"

var gStreamListInfo = TRNStreamListInfoST{
	Streams: make(TRNStreamsMAP),
	//Streams_extra: make(TRNStreamsMAP),
	pseudoUUID: func() (uuid string) {
		b := make([]byte, 16)
		_, err := rand.Read(b)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		uuid = fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
		return
	},
}

type TRNStreamListInfoST struct {
	mutex   sync.RWMutex
	Streams TRNStreamsMAP `json:"streams" groups:"config"`
	//Streams_extra TRNStreamsMAP `json:"streams_extra" groups:"config"`
	//LastError error
	pseudoUUID func() (uuid string)
}

type TRNStreamsMAP map[string]TRNStreamST
type TRNStreamST struct {
	//StreamId string `json:"stream_id" groups:"config"`
	//CctvName string `json:"cctv_name" groups:"config"`
	//CctvIp   string `json:"cctv_ip" groups:"config" `
	////Channels     ChannelMAP
	RtspUrl string `json:"source" groups:"config" binding:"required"`
	//RtspUrl_2 string `json:"url2" groups:"config"`
	//Status    byte   `json:"status" groups:"config"` //healty
	OnDemand bool `json:"sourceOnDemand" groups:"config"`
	////DisableAudio bool   `json:"disable_audio" groups:"config"`
	////Debug        bool   `json:"debug" groups:"config"`
}

// func (obj *TRNStreamListInfoST) apply_to_list(newStreamsList StreamsMAP) bool {
// 	obj.mutex.Lock()
// 	defer obj.mutex.Unlock()

// 	isListChanged := false
// 	var streamToDelete []string //group of iSuuid

// 	//same suuid -> check
// 	for iSuuid, oldStream := range obj.Streams {
// 		if newStream, ok := (newStreamsList)[iSuuid]; ok {
// 			change_found := false
// 			if oldStream.RtspUrl != newStream.RtspUrl { //different -> change
// 				change_found = true
// 				oldStream.RtspUrl = newStream.RtspUrl
// 			}
// 			// if oldStream.RtspUrl_2 != newStream.RtspUrl_2 { //different -> change
// 			// 	change_found = true
// 			// 	oldStream.RtspUrl_2 = newStream.RtspUrl_2
// 			// }
// 			if oldStream.Status != newStream.Status { //different -> change
// 				change_found = true
// 				oldStream.Status = newStream.Status
// 			}
// 			if oldStream.CctvName != newStream.CctvName { //different -> change
// 				change_found = true
// 				oldStream.CctvName = newStream.CctvName
// 			}

// 			if change_found {
// 				(obj.Streams)[iSuuid] = oldStream
// 				if !isListChanged {
// 					isListChanged = true
// 				}
// 			}
// 			delete(newStreamsList, iSuuid)
// 		} else {
// 			streamToDelete = append(streamToDelete, iSuuid)
// 		}
// 	}

// 	//no suuid ->delete
// 	for _, iSuuid := range streamToDelete {
// 		delete((obj.Streams), iSuuid)
// 		if !isListChanged {
// 			isListChanged = true
// 		}
// 	}

// 	//new suuid ->add
// 	for iSuuid, newStream := range newStreamsList {
// 		(obj.Streams)[iSuuid] = newStream
// 		if !isListChanged {
// 			isListChanged = true
// 		}
// 	}

// 	return isListChanged

// }

func (obj *TRNStreamListInfoST) SaveList() error {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	// log.WithFields(logrus.Fields{
	// 	"module": "stream_list",
	// 	"func":   "NewStreamCore",
	// }).Debugln("Saving configuration to", StreamListJsonFile)
	v2, err := goversion.NewVersion("2.0.0")
	if err != nil {
		return err
	}

	options := &sheriff.Options{
		Groups:     []string{"config"},
		ApiVersion: v2,
	}
	data, err := sheriff.Marshal(options, obj)
	if err != nil {
		return err
	}
	//data := obj
	JsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(STREAM_LIST_FILENAME, JsonData, 0644)
	if err != nil {
		// log.WithFields(logrus.Fields{
		// 	"module": "stream_list",
		// 	"func":   "SaveList",
		// 	"call":   "WriteFile",
		// }).Errorln(err.Error())
		return err
	}

	return nil
}
