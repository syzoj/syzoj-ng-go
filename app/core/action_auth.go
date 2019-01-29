package core

import (
    "errors"
    "context"
    "regexp"
    "time"

    "github.com/mongodb/mongo-go-driver/mongo"
    mongo_options "github.com/mongodb/mongo-go-driver/mongo/options"
    "github.com/mongodb/mongo-go-driver/bson"
    "github.com/mongodb/mongo-go-driver/bson/primitive"
    "golang.org/x/crypto/bcrypt"
)

var userNamePattern = regexp.MustCompile("^[0-9A-Za-z]{3,32}$")

func checkUserName(userName string) bool {
	return userNamePattern.MatchString(userName)
}

type Register1 struct {
    UserName string
    Password string
}
type Register1Resp struct {
    UserId primitive.ObjectID
}

var ErrConflict = errors.New("Conflict operation")
var ErrDuplicateUserName = errors.New("Duplicate user name")
var ErrInvalidUserName = errors.New("Invalid user name")

// Possible errors:
// * nil: success
// * ErrInvalidUserName
// * ErrDuplicateUserName
// * Other MongoDB errors or context errors
func (c *Core) Action_Register(ctx context.Context, req *Register1) (*Register1Resp, error) {
    var err error
    if !checkUserName(req.UserName) {
        return nil, ErrInvalidUserName
    }
    var passwordHash []byte
    if passwordHash, err = bcrypt.GenerateFromPassword([]byte(req.Password), 0); err != nil {
        panic(err)
    }
    lock := c.LockOracle([]interface{}{KeyUserName(req.UserName)})
    if lock == nil {
        return nil, ErrConflict
    }
    defer lock.Release()
    if _, err = c.mongodb.Collection("user").FindOne(ctx, bson.D{{"username", req.UserName}}, mongo_options.FindOne().SetProjection(bson.D{{"_id", 1}})).DecodeBytes(); err != nil {
        if err != mongo.ErrNoDocuments {
            return nil, err
        }
    } else {
        return nil, ErrDuplicateUserName
    }
    userId := primitive.NewObjectID()
    if _, err = c.mongodb.Collection("user").InsertOne(ctx, bson.D{
        {"_id", userId},
        {"username", req.UserName},
        {"register_time", time.Now()},
        {"auth", bson.D{{"password", passwordHash}, {"method", int64(1)}}},
    }); err != nil {
        return nil, err
    }
    log.WithField("username", req.UserName).Info("Created account")
    return &Register1Resp{
        UserId: userId,
    }, nil
}
