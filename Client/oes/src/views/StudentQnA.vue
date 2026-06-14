<template>
  <div class="main-content">
    <section class="section-wrapper">
      <!-- Section title -->
      <div class="row section-title">
        <div class="col-12">
          <h4>Ask a Question about Course Material</h4>
          <p class="text-muted">
            Powered by a local AI model — your questions never leave this server.
          </p>
        </div>
      </div>

      <!-- Chat window -->
      <div class="row section-body">
        <div class="col-12">
          <div class="chat-window" ref="chatWindow">
            <div
              v-for="(msg, index) in messages"
              :key="index"
              :class="['chat-bubble', msg.role === 'student' ? 'bubble-student' : 'bubble-ai']"
            >
              <span class="bubble-label">{{ msg.role === 'student' ? 'You' : 'AI Assistant' }}</span>
              <p class="bubble-text">{{ msg.text }}</p>
            </div>

            <!-- Typing indicator -->
            <div v-if="loading" class="chat-bubble bubble-ai">
              <span class="bubble-label">AI Assistant</span>
              <p class="bubble-text typing-indicator">Thinking<span>.</span><span>.</span><span>.</span></p>
            </div>
          </div>
        </div>
      </div>

      <!-- Optional: context/topic filter -->
      <div class="row mt-2">
        <div class="col-md-4">
          <div class="form-group">
            <label for="contextId">Topic / Subject (optional)</label>
            <input
              id="contextId"
              v-model="contextId"
              type="text"
              class="form-control"
              placeholder="e.g. biology, chapter-3"
              :disabled="loading"
            />
          </div>
        </div>
      </div>

      <!-- Question input -->
      <div class="row mt-2">
        <div class="col-12">
          <div class="input-group">
            <input
              v-model="question"
              type="text"
              class="form-control"
              placeholder="Type your question here..."
              :disabled="loading"
              @keyup.enter="submitQuestion"
            />
            <div class="input-group-append">
              <button
                class="btn btn-primary"
                @click="submitQuestion"
                :disabled="loading || !question.trim()"
              >
                <span v-if="loading">
                  <span class="spinner-border spinner-border-sm" role="status"></span>
                  Asking...
                </span>
                <span v-else>Ask</span>
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Error display -->
      <div v-if="errorMsg" class="row mt-2">
        <div class="col-12">
          <div class="alert alert-danger">{{ errorMsg }}</div>
        </div>
      </div>
    </section>
  </div>
</template>

<script>
export default {
  name: "StudentQnA",
  data() {
    return {
      question: "",
      contextId: "",
      messages: [],
      loading: false,
      errorMsg: "",
    };
  },
  methods: {
    async submitQuestion() {
      const q = this.question.trim();
      if (!q || this.loading) return;

      // Add student message to chat
      this.messages.push({ role: "student", text: q });
      this.question = "";
      this.errorMsg = "";
      this.loading = true;

      // Scroll to bottom
      this.$nextTick(() => this.scrollToBottom());

      try {
        const payload = { question: q };
        if (this.contextId.trim()) {
          payload.contextId = this.contextId.trim();
        }

        const res = await this.$http.post("/api/r/askQuestion", payload);

        const answer = res.data.answer || "No answer returned.";
        this.messages.push({ role: "ai", text: answer });
      } catch (err) {
        const msg =
          err.response && err.response.data && err.response.data.error
            ? err.response.data.error
            : "Failed to get answer. Please try again.";
        this.errorMsg = msg;
        this.messages.push({ role: "ai", text: `Error: ${msg}` });
      } finally {
        this.loading = false;
        this.$nextTick(() => this.scrollToBottom());
      }
    },

    scrollToBottom() {
      const el = this.$refs.chatWindow;
      if (el) {
        el.scrollTop = el.scrollHeight;
      }
    },
  },
};
</script>

<style lang="scss" scoped>
.chat-window {
  border: 1px solid #dee2e6;
  border-radius: 8px;
  padding: 16px;
  min-height: 300px;
  max-height: 500px;
  overflow-y: auto;
  background: #f8f9fa;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.chat-bubble {
  max-width: 75%;
  padding: 10px 14px;
  border-radius: 12px;
  word-break: break-word;
}

.bubble-student {
  align-self: flex-end;
  background: #007bff;
  color: #fff;
  border-bottom-right-radius: 2px;
}

.bubble-ai {
  align-self: flex-start;
  background: #ffffff;
  border: 1px solid #dee2e6;
  color: #212529;
  border-bottom-left-radius: 2px;
}

.bubble-label {
  display: block;
  font-size: 0.7rem;
  font-weight: 600;
  margin-bottom: 4px;
  opacity: 0.7;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.bubble-text {
  margin: 0;
  white-space: pre-wrap;
  line-height: 1.5;
}

.typing-indicator span {
  animation: blink 1.2s infinite;
  &:nth-child(2) { animation-delay: 0.2s; }
  &:nth-child(3) { animation-delay: 0.4s; }
}

@keyframes blink {
  0%, 80%, 100% { opacity: 0; }
  40% { opacity: 1; }
}
</style>
