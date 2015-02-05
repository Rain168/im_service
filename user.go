package main
import "math/rand"
import "fmt"
import log "github.com/golang/glog"
import "github.com/garyburd/redigo/redis"


const CHARACTER_SET = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenUserToken() string {
	b := make([]byte, 30)
	for i := 0; i < 30; i++ {
		r := rand.Int()%len(CHARACTER_SET)
		b[i] = CHARACTER_SET[r]
	}
	return string(b)
}

func GetUserAccessToken(appid int64, uid int64) string {
	conn := redis_pool.Get()
	defer conn.Close()

	key := fmt.Sprintf("users_%d_%d", appid, uid)
	token, err := redis.String(conn.Do("HGET", key, "access_token"))
	if err != nil {
		log.Infof("hget %s err:%s\n", key, err)
		return ""
	}
	return token
}

func LoadUserAccessToken(token string) (int64, int64, string, error) {
	conn := redis_pool.Get()
	defer conn.Close()

	key := fmt.Sprintf("access_token_%s", token)

	var uid int64
	var appid int64
	var uname string
	reply, err := redis.Values(conn.Do("HMGET", key, "user_id", "app_id", "user_name"))
	if err != nil {
		log.Info("hmget error:", err)
		return 0, 0, "", err
	}

	_, err = redis.Scan(reply, &uid, &appid, &uname)
	if err != nil {
		log.Warning("scan error:", err)
		return 0, 0, "", err
	}
	return appid, uid, uname, nil	
}

func SaveUserAccessToken(appid int64, uid int64, uname string, token string) error {
	conn := redis_pool.Get()
	defer conn.Close()

	key := fmt.Sprintf("access_token_%s", token)
	
	_, err := conn.Do("HMSET", key, "user_id", uid, "user_name", uname, "app_id", appid)
	if err != nil {
		log.Info("hmset err:", err)
		return err
	}

	key = fmt.Sprintf("users_%d_%d", appid, uid)
	_, err = conn.Do("HSET", key, "access_token", token)
	if err != nil {
		log.Info("hget err:", err)
		return err
	}
	return nil
}

func SaveUserDeviceToken(appid int64, uid int64, device_token string, ng_device_token string) error {
	conn := redis_pool.Get()
	defer conn.Close()

	key := fmt.Sprintf("users_%d_%d", appid, uid)
	_, err := conn.Do("HMSET", key, "apns_device_token", device_token, 
		"ng_device_token", ng_device_token)
	if err != nil {
		log.Info("hget err:", err)
		return err
	}
	return nil	
}
