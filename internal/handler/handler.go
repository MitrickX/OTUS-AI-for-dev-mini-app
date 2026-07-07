package handler

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/dmitrypavlov/mini-questionnaire/api"
)

var questions = []api.Question{
	{
		Id:   uuid.MustParse("a1b2c3d4-e5f6-7890-abcd-ef1234567890"),
		Text: "Как вас зовут?",
		Type: api.Text,
		Required: boolPtr(true),
	},
	{
		Id:   uuid.MustParse("b2c3d4e5-f6a7-8901-bcde-f12345678901"),
		Text: "Какой ваш любимый цвет?",
		Type: api.SingleChoice,
		Options: &[]string{"Красный", "Синий", "Зелёный", "Жёлтый", "Другой"},
		Required: boolPtr(true),
	},
	{
		Id:   uuid.MustParse("c3d4e5f6-a7b8-9012-cdef-123456789012"),
		Text: "Какими языками программирования вы владеете?",
		Type: api.MultipleChoice,
		Options: &[]string{"Go", "Python", "JavaScript", "Java", "C++", "Rust", "Другой"},
	},
	{
		Id:   uuid.MustParse("d4e5f6a7-b8c9-0123-defa-234567890123"),
		Text: "Сколько лет вы занимаетесь программированием?",
		Type: api.SingleChoice,
		Options: &[]string{"Меньше года", "1–3 года", "3–5 лет", "5–10 лет", "Больше 10 лет"},
		Required: boolPtr(true),
	},
	{
		Id:   uuid.MustParse("e5f6a7b8-c9d0-1234-efab-345678901234"),
		Text: "Что бы вы хотели улучшить в нашем продукте?",
		Type: api.Text,
	},
}

type Server struct {
	mu       sync.Mutex
	answers  []api.AnswerRecord
}

func New() *Server {
	return &Server{}
}

func (s *Server) GetQuestions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(questions)
}

func (s *Server) SubmitAnswers(w http.ResponseWriter, r *http.Request) {
	var req api.SubmitAnswersRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(api.Error{Error: "invalid JSON body"})
		return
	}

	if len(req.Answers) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(api.Error{Error: "answers must not be empty"})
		return
	}

	record := api.AnswerRecord{
		Id:          uuid.New(),
		Respondent:  req.Respondent,
		Answers:     req.Answers,
		SubmittedAt: time.Now().UTC(),
	}

	s.mu.Lock()
	s.answers = append(s.answers, record)
	s.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(record)
}

func boolPtr(v bool) *bool {
	return &v
}
