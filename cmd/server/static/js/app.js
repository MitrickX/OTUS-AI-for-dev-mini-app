const API_BASE = '';

let questions = [];

async function loadQuestions() {
  const el = document.getElementById('questions-container');
  el.innerHTML = '<div class="loading">Загрузка вопросов...</div>';

  try {
    const res = await fetch(API_BASE + '/questions');
    if (!res.ok) throw new Error('Ошибка загрузки вопросов');
    questions = await res.json();
    renderQuestions();
  } catch (err) {
    el.innerHTML = '<div class="error">Не удалось загрузить вопросы. Попробуйте позже.</div>';
  }
}

function renderQuestions() {
  const el = document.getElementById('questions-container');
  el.innerHTML = '';

  questions.forEach((q, idx) => {
    const div = document.createElement('div');
    div.className = 'question';

    const label = document.createElement('span');
    label.className = 'question-label';
    label.textContent = `${idx + 1}. ${q.text}`;
    if (q.required) {
      const star = document.createElement('span');
      star.className = 'required';
      star.textContent = '*';
      label.appendChild(star);
    }
    div.appendChild(label);

    switch (q.type) {
      case 'text':
        const input = document.createElement('input');
        input.type = 'text';
        input.dataset.qid = q.id;
        input.placeholder = 'Ваш ответ...';
        div.appendChild(input);
        break;

      case 'single_choice':
        (q.options || []).forEach(opt => {
          const wrap = document.createElement('div');
          wrap.className = 'option';
          const radio = document.createElement('input');
          radio.type = 'radio';
          radio.name = `q_${q.id}`;
          radio.value = opt;
          radio.dataset.qid = q.id;
          const lbl = document.createElement('label');
          lbl.textContent = opt;
          wrap.appendChild(radio);
          wrap.appendChild(lbl);
          div.appendChild(wrap);
        });
        break;

      case 'multiple_choice':
        (q.options || []).forEach(opt => {
          const wrap = document.createElement('div');
          wrap.className = 'option';
          const cb = document.createElement('input');
          cb.type = 'checkbox';
          cb.value = opt;
          cb.dataset.qid = q.id;
          const lbl = document.createElement('label');
          lbl.textContent = opt;
          wrap.appendChild(cb);
          wrap.appendChild(lbl);
          div.appendChild(wrap);
        });
        break;
    }

    el.appendChild(div);
  });

  document.getElementById('submit-area').classList.remove('hidden');
}

function collectAnswers() {
  const answers = [];

  questions.forEach(q => {
    const els = document.querySelectorAll(`[data-qid="${q.id}"]`);

    switch (q.type) {
      case 'text': {
        const val = els[0]?.value?.trim();
        if (val) answers.push({ question_id: q.id, value: val });
        break;
      }
      case 'single_choice': {
        const checked = document.querySelector(`input[name="q_${q.id}"]:checked`);
        if (checked) answers.push({ question_id: q.id, value: checked.value });
        break;
      }
      case 'multiple_choice': {
        const checked = [...els].filter(el => el.checked).map(el => el.value);
        if (checked.length > 0) answers.push({ question_id: q.id, value: checked });
        break;
      }
    }
  });

  return answers;
}

async function submitAnswers() {
  const msgEl = document.getElementById('message');
  msgEl.classList.add('hidden');
  msgEl.textContent = '';

  const answers = collectAnswers();

  if (answers.length === 0) {
    msgEl.className = 'error';
    msgEl.textContent = 'Заполните хотя бы один вопрос.';
    msgEl.classList.remove('hidden');
    return;
  }

  const btn = document.querySelector('.submit-btn');
  btn.disabled = true;
  btn.textContent = 'Отправка...';

  try {
    const res = await fetch(API_BASE + '/answers', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ answers }),
    });

    if (res.ok) {
      msgEl.className = 'success';
      msgEl.textContent = 'Спасибо! Ваши ответы сохранены.';
    } else {
      const err = await res.json();
      msgEl.className = 'error';
      msgEl.textContent = err.error || 'Ошибка при отправке.';
    }
  } catch {
    msgEl.className = 'error';
    msgEl.textContent = 'Ошибка соединения с сервером.';
  }

  msgEl.classList.remove('hidden');
  btn.disabled = false;
  btn.textContent = 'Отправить';
}

document.addEventListener('DOMContentLoaded', loadQuestions);
