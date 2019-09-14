package judgecore

import (
	"strconv"
	"sync"
	"time"

	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go"
)

// The Server object contains all information to run a server.
// The expected usage is to fill in all exported fields and then call Serve().
// No mutations to the struct should be made after calling Serve().
type Server struct {
	ListenAddr string        // The tcp address to listen on.
	Minio      *minio.Client // The MinIO client for S3 backend.
	FileBucket string        // The name of the bucket to store uploaded files.

	mu sync.Mutex
}

// The request body for POST /task/create.
type CreateTaskRequest struct {
	FileObjectName string `json:"file_object_name"` // The name of the uploaded file in S3 bucket.
}

// POST /task/create
func (s *Server) PostTask(c *gin.Context) {
	var req CreateTaskRequest
	if c.Bind(&req) != nil {
		return
	}
	c.JSON(200, req)
}

// GET /judger/listen
// A judger makes a request to this endpoint to listen for tasks, and sends results to POST /judger/post.
// The server sends judge messages in the form of SSE (Server Sent Events).
// The session is as long as the connection. Once the connection drops, the judger is considered
// offline and the tasks are either marked as failed or requeued.
func (s *Server) GetJudgerListen(c *gin.Context) {
	var id int
	// Send the first event with SSEvent to properly set Content-Type.
	// Alternatively call sse.Event.WriteContentType on c.Writer
	c.Render(-1, sse.Event{
		Event: "open",
		Data:  "",
	})
	log.Info("listen")
loop:
	for {
		select {
		case <-c.Request.Context().Done():
			log.Info("request cancelled")
			break loop
		case <-time.After(time.Second):
			log.Info("event")
			id++
			c.Render(-1, sse.Event{
				Id:    strconv.Itoa(id),
				Event: "message",
				Data:  "data",
			})
			c.Writer.Flush()
		}
	}
}

// Starts the server.
func (s *Server) Serve() error {
	router := gin.Default()
	router.POST("/task/create", s.PostTask)
	router.GET("/judger/listen", s.GetJudgerListen)
	return router.Run(s.ListenAddr)
}
