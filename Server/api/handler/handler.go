package handler

import (
	"sync"

	"github.com/mskKandula/oes/api/model"
)

// QuestionCache stores exam questions temporarily in memory with thread-safe access
type QuestionCache struct {
	mu    sync.RWMutex
	cache map[int64][]string // examId -> questions
}

func NewQuestionCache() *QuestionCache {
	return &QuestionCache{
		cache: make(map[int64][]string),
	}
}

func (qc *QuestionCache) Set(examId int64, questions []string) {
	qc.mu.Lock()
	defer qc.mu.Unlock()
	qc.cache[examId] = questions
}

func (qc *QuestionCache) Get(examId int64) ([]string, bool) {
	qc.mu.RLock()
	defer qc.mu.RUnlock()
	questions, exists := qc.cache[examId]
	return questions, exists
}

type Handler struct {
	UserService    model.UserService
	StudentService model.StudentService
	CommonService  model.CommonService
	QuestionCache  *QuestionCache
}

func NewHandler(us model.UserService, ss model.StudentService, cs model.CommonService) *Handler {
	return &Handler{
		UserService:    us,
		StudentService: ss,
		CommonService:  cs,
		QuestionCache:  NewQuestionCache(),
	}
}
