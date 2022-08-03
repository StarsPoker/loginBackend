package one_time_password

import (
	"github.com/StarsPoker/loginBackend/logger"

	"github.com/StarsPoker/loginBackend/datasources/mongo/stars_mongo"
	"github.com/StarsPoker/loginBackend/domain/access_token"
	"github.com/StarsPoker/loginBackend/utils/errors/rest_errors"

	"github.com/StarsPoker/loginBackend/utils/date_utils"
	"github.com/StarsPoker/loginBackend/utils/mongo_utils"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	database   = "stars"
	collection = "one_time_password"
)

func Insert(otp OneTimePassword) *rest_errors.RestErr {
	session, _ := stars_mongo.GetSession()
	defer session.Close()

	col := session.DB(database).C(collection)
	err := col.Insert(&otp)

	if err != nil {
		logger.Error("error when trying to creat a one_time_password", err)
		return rest_errors.NewInternalServerError("database error")
	}

	return nil
}

func GetAuth(accessTokenRequest access_token.AccessTokenRequest) (*OneTimePassword, *rest_errors.RestErr) {

	otp := OneTimePassword{}
	oneTimePasswordKey := accessTokenRequest.ClientKey
	userId := accessTokenRequest.ClientId

	session, _ := stars_mongo.GetSession()
	defer session.Close()

	col := session.DB(database).C(collection)
	err := col.Find(bson.M{"key": oneTimePasswordKey, "user_id": userId}).One(&otp)

	if err != nil {
		logger.Error("error when trying to get a authentication", err)
		return nil, mongo_utils.ParseError(err)
	}
	return &otp, nil
}

func DeleteAllById(otp OneTimePassword, rec bool) bool {
	if rec {
		return true
	} else {
		session, _ := stars_mongo.GetSession()
		defer session.Close()

		col := session.DB(database).C(collection)
		err := col.Remove(bson.M{"user_id": otp.UserId})

		if err == nil {
			return DeleteAllById(otp, false)
		}
		return true
	}

}

func DeleteByKey(otp OneTimePassword) *rest_errors.RestErr {
	session, _ := stars_mongo.GetSession()
	defer session.Close()

	col := session.DB(database).C(collection)
	err := col.Remove(bson.M{"key": otp.Key})

	if err != nil {
		logger.Error("error when trying to remove a one_time_password", err)
		return rest_errors.NewInternalServerError("database error")
	}

	return nil
}

func DeleteExpiredOneTimePasswords() *rest_errors.RestErr {
	session, _ := stars_mongo.GetSession()
	defer session.Close()

	col := session.DB(database).C(collection)

	err := col.Remove(bson.M{"expires": bson.M{"$lte": date_utils.GetNow().Unix()}})

	if err != nil {
		logger.Error("error when trying to delete a one_time_password", err)
		return rest_errors.NewInternalServerError("database error")
	}
	return nil
}
