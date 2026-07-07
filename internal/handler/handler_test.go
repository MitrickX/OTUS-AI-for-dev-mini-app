package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dmitrypavlov/mini-questionnaire/api"
)

func setupTest() *httptest.Server {
	s := New()
	mux := http.NewServeMux()
	mux.HandleFunc("GET /questions", s.GetQuestions)
	mux.HandleFunc("POST /answers", s.SubmitAnswers)
	return httptest.NewServer(mux)
}

func TestGetQuestions_Returns200(t *testing.T) {
	ts := setupTest()
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/questions")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var questions []api.Question
	if err := json.NewDecoder(resp.Body).Decode(&questions); err != nil {
		t.Fatal(err)
	}

	if len(questions) != 5 {
		t.Errorf("expected 5 questions, got %d", len(questions))
	}
}

func TestGetQuestions_CheckTypes(t *testing.T) {
	ts := setupTest()
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/questions")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	var questions []api.Question
	json.NewDecoder(resp.Body).Decode(&questions)

	types := map[int]api.QuestionType{
		0: api.Text,
		1: api.SingleChoice,
		2: api.MultipleChoice,
		3: api.SingleChoice,
		4: api.Text,
	}

	for i, q := range questions {
		if q.Type != types[i] {
			t.Errorf("question[%d] expected type %s, got %s", i, types[i], q.Type)
		}
	}
}

func TestGetQuestions_ContentType(t *testing.T) {
	ts := setupTest()
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/questions")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if ct := resp.Header.Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", ct)
	}
}

func TestSubmitAnswers_ValidRequest(t *testing.T) {
	ts := setupTest()
	defer ts.Close()

	body := api.SubmitAnswersRequest{
		Respondent: strPtr("Иван"),
		Answers: []api.Answer{
			{QuestionId: questions[0].Id, Value: mustAnswerValue("Test")},
		},
	}

	b, _ := json.Marshal(body)
	resp, err := http.Post(ts.URL+"/answers", "application/json", bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected status 201, got %d", resp.StatusCode)
	}

	var record api.AnswerRecord
	if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
		t.Fatal(err)
	}

	if record.Id.String() == "00000000-0000-0000-0000-000000000000" {
		t.Error("expected non-zero UUID")
	}

	if record.Respondent == nil || *record.Respondent != "Иван" {
		t.Errorf("expected respondent Иван, got %v", record.Respondent)
	}

	if record.SubmittedAt.IsZero() {
		t.Error("expected non-zero submitted_at")
	}

	if len(record.Answers) != 1 {
		t.Errorf("expected 1 answer, got %d", len(record.Answers))
	}
}

func TestSubmitAnswers_EmptyBody(t *testing.T) {
	ts := setupTest()
	defer ts.Close()

	resp, err := http.Post(ts.URL+"/answers", "application/json", bytes.NewReader([]byte{}))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestSubmitAnswers_EmptyAnswers(t *testing.T) {
	ts := setupTest()
	defer ts.Close()

	body := api.SubmitAnswersRequest{
		Answers: []api.Answer{},
	}
	b, _ := json.Marshal(body)
	resp, err := http.Post(ts.URL+"/answers", "application/json", bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}

	var errResp api.Error
	json.NewDecoder(resp.Body).Decode(&errResp)
	if errResp.Error != "answers must not be empty" {
		t.Errorf("unexpected error message: %s", errResp.Error)
	}
}

func TestSubmitAnswers_MultipleSubmissions(t *testing.T) {
	ts := setupTest()
	defer ts.Close()

	body := api.SubmitAnswersRequest{
		Answers: []api.Answer{
			{QuestionId: questions[0].Id, Value: mustAnswerValue("A")},
		},
	}
	b, _ := json.Marshal(body)

	for i := 0; i < 3; i++ {
		resp, err := http.Post(ts.URL+"/answers", "application/json", bytes.NewReader(b))
		if err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()
	}

	s := New()
	if len(s.answers) != 0 {
		t.Errorf("fresh server should have 0 answers, got %d", len(s.answers))
	}
}

func TestSubmitAnswers_ContentType(t *testing.T) {
	ts := setupTest()
	defer ts.Close()

	body := api.SubmitAnswersRequest{
		Answers: []api.Answer{
			{QuestionId: questions[0].Id, Value: mustAnswerValue("test")},
		},
	}
	b, _ := json.Marshal(body)
	resp, err := http.Post(ts.URL+"/answers", "application/json", bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if ct := resp.Header.Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", ct)
	}
}

func strPtr(s string) *string {
	return &s
}

func mustAnswerValue(v string) api.Answer_Value {
	var av api.Answer_Value
	if err := av.FromAnswerValue0(v); err != nil {
		panic(err)
	}
	return av
}
