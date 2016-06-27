package router

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"encoding/json"
	"net/http"
	"github.com/ugorji/go/codec"
	"strings"
	"github.com/ruslanfirefly/bolt-wrap"
	"restdb/common_utils"
)

func GetRouter(dbpath string) *gin.Engine {
	r := gin.Default()
	dataMap := make(map[string]interface{})
	db := bolt_wrap.New(dbpath)

	r.POST("db/:bucket/:data", func(c *gin.Context) {
		bucket, _ := c.Params.Get("bucket")
		key, _ := c.Params.Get("data")
		data, _ := ioutil.ReadAll(c.Request.Body)
		if tenp := db.Get(bucket, key); string(tenp) == "" {
			err := json.Unmarshal(data, &dataMap)
			if err != nil {
				c.String(http.StatusBadRequest, "Can not parse JSON")
			} else {
				var dataMap_bytes  []byte
				codec.NewEncoderBytes(&dataMap_bytes, new(codec.CborHandle)).Encode(dataMap)
				err := db.Set(bucket, key, dataMap_bytes)
				if err != nil {
					c.String(http.StatusBadRequest, "Can not save data")
				} else {
					c.String(http.StatusOK, "ok")
				}
			}
		} else {
			c.String(http.StatusConflict, "This key exist")
		}

	})

	r.GET("db/:bucket/:data/:fields", func(c *gin.Context) {
		bucket, _ := c.Params.Get("bucket")

		key, _ := c.Params.Get("data")
		fields, _ := c.Params.Get("fields")
		data_byte := db.Get(bucket, key)
		data := make(map[interface{}]interface{})
		resData := make(map[string]interface{})
		codec.NewDecoderBytes(data_byte, new(codec.CborHandle)).Decode(&data)
		resData = common_utils.ParseMap(data)
		if fields == "all" {
			c.JSON(http.StatusOK, resData)
		} else {
			arrFields := strings.Split(fields, "&")
			resData2 := make(map[string]interface{})
			for _, v := range arrFields {
				resData2[v] = resData[v]
			}
			c.JSON(http.StatusOK, resData2)
		}

	})

	r.PUT("db/:bucket/:data", func(c *gin.Context) {
		bucket, _ := c.Params.Get("bucket")
		key, _ := c.Params.Get("data")
		data, _ := ioutil.ReadAll(c.Request.Body)
		dataFromBase := db.Get(bucket, key);
		dataMapFromReq := make(map[string]interface{})
		if string(dataFromBase) != "" {
			err := json.Unmarshal(data, &dataMapFromReq)
			if err != nil {
				c.String(http.StatusBadRequest, "Can not parse JSON")
			} else {
				dataMap := make(map[interface{}]interface{})
				codec.NewDecoderBytes(dataFromBase, new(codec.CborHandle)).Decode(&dataMap)
				resData := common_utils.ParseMap(dataMap)
				for key := range dataMapFromReq {
					if str, ok := dataMapFromReq[key].(string); ok && str == "" {
						delete(resData, key)
					} else {
						resData[key] = dataMapFromReq[key]
					}

				}
				var dataMap_bytes  []byte
				codec.NewEncoderBytes(&dataMap_bytes, new(codec.CborHandle)).Encode(resData)
				err := db.Set(bucket, key, dataMap_bytes)
				if err != nil {
					c.String(http.StatusBadRequest, "Can not update data")
				} else {
					c.String(http.StatusOK, "ok")
				}
			}
		} else {
			c.String(http.StatusConflict, "This key not exist")
		}
	})

	r.DELETE("db/:bucket/:data", func(c *gin.Context) {
		bucket, _ := c.Params.Get("bucket")
		key, _ := c.Params.Get("data")
		err := db.Delete(bucket, key)
		if err == nil {
			c.String(http.StatusOK, "ok")
		} else {
			c.String(http.StatusBadRequest, "User not exist")
		}
	})

	r.GET("/backup/:filename", func(c *gin.Context) {
		filename, _ := c.Params.Get("filename")
		db.GetBackUp(c.Writer, filename)
	})

	return r
}
