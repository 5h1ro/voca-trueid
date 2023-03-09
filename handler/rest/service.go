package rest

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"
	"vocatrueid/entity"
	"vocatrueid/helpers"
	"vocatrueid/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"github.com/thedevsaddam/govalidator"
)

type Service struct {
	utilsService utils.Service
}

func NewService(utilsService utils.Service) *Service {
	return &Service{utilsService}
}

func (s *Service) CheckNickname(c *gin.Context) {
	var body entity.User
	region := c.Param("region")
	_ = region
	code := c.Param("code")

	rules := govalidator.MapData{
		"target": []string{"required"},
	}

	opts := govalidator.Options{
		Request: c.Request,
		Data:    &body,
		Rules:   rules,
	}

	v := govalidator.New(opts)
	validate := v.ValidateJSON()
	err := map[string]interface{}{"validationError": validate}
	if len(validate) != 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  err,
		})
		return
	}

	requestLog, _ := OpenLogFile("logs/request.log")
	log.SetOutput(requestLog)
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
	log.Println("IP: " + c.ClientIP() + ", target: " + body.Target + ", secret: " + body.Secret)

	if body.Secret == "" || helpers.GetEnv("SECRET_KEY") != body.Secret {
		res := map[string]interface{}{
			"message": "Invalid Secret",
			"data":    map[string]interface{}{},
		}
		c.JSON(http.StatusUnauthorized, res)
		return
	}

	ring := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			"server": ":6379",
		},
	})
	mycache := cache.New(&cache.Options{
		Redis:      ring,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})
	ctx := context.TODO()
	key := body.Target

	var wanted entity.Object
	if err := mycache.Get(ctx, key, &wanted); err == nil {
		res := map[string]interface{}{
			"message": "Success check username",
			"data": map[string]interface{}{
				"username": wanted.Username,
			},
		}
		c.JSON(http.StatusOK, res)
		return
	}

	resultLog, _ := OpenLogFile("logs/result.log")
	log.SetOutput(resultLog)
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
	CheckNickname := s.utilsService.CheckNickname(body, code)
	if CheckNickname.IsSuccess {
		res := map[string]interface{}{
			"message": "Success check username",
			"data": map[string]interface{}{
				"username": CheckNickname.Username,
			},
		}
		obj := &entity.Object{
			Username: CheckNickname.Username,
		}

		if err := mycache.Set(&cache.Item{
			Ctx:   ctx,
			Key:   key,
			Value: obj,
			TTL:   72 * time.Hour,
		}); err != nil {
			panic(err)
		}
		log.Println("Target: " + body.Target + ", Username: " + CheckNickname.Username)
		c.JSON(http.StatusOK, res)
	} else {
		if CheckNickname.Username != "" {
			res := map[string]interface{}{
				"message": CheckNickname.Username,
				"data":    map[string]interface{}{},
			}
			c.JSON(http.StatusForbidden, res)
		} else {
			res := map[string]interface{}{
				"message": "Invalid User",
				"data":    map[string]interface{}{},
			}
			c.JSON(http.StatusNotFound, res)
		}
	}
}

func OpenLogFile(path string) (*os.File, error) {
	logFile, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	return logFile, nil
}
