"use client";

import { useState, useRef, useEffect, useCallback } from "react";
import { AnimatePresence } from "framer-motion";
import { useChatMessages } from "@/hooks/useChatMessages";
import { useVoiceRecording } from "@/hooks/useVoiceRecording";
import { EditorEmptyState } from "@/components/editor/EditorEmptyState";
import { ChatMessage } from "@/components/editor/ChatMessage";
import { TypingIndicator } from "@/components/editor/TypingIndicator";
import { ChatInput } from "@/components/editor/ChatInput";
import { SqlEditorSheet } from "@/components/editor/SqlEditorSheet";

const suggestedQueries = [
  "Show me sales trends for last quarter",
  "What are the top performing products?",
  "Analyze customer demographics",
  "Compare revenue across regions",
];

export default function EditorPage() {
  const [input, setInput] = useState("");
  const [showSqlEditor, setShowSqlEditor] = useState(false);
  const [sqlQuery, setSqlQuery] = useState("");
  const messagesEndRef = useRef<HTMLDivElement>(null);

  const {
    messages,
    isSubmitting,
    audioProcessing,
    copied,
    handleSend,
    handleSqlSubmit,
    processAudio,
    copyToClipboard,
    clearChat,
  } = useChatMessages();

  const {
    isRecording,
    recordedBlob,
    startRecording,
    stopRecording,
    discardRecording,
  } = useVoiceRecording();

  const scrollToBottom = useCallback(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  }, []);

  useEffect(() => {
    scrollToBottom();
  }, [messages, scrollToBottom]);

  const onSend = () => {
    if (!input.trim()) return;
    handleSend(input);
    setInput("");
  };

  const onSqlSubmit = () => {
    if (!sqlQuery.trim()) return;
    handleSqlSubmit(sqlQuery);
    setShowSqlEditor(false);
    setSqlQuery("");
  };

  const onProcessAudio = () => {
    if (recordedBlob) {
      processAudio(recordedBlob);
      discardRecording();
    }
  };

  return (
    <div className="flex flex-col h-full bg-background">
      {/* Messages area */}
      <div className="flex-1 overflow-y-auto">
        <AnimatePresence mode="wait">
          {messages.length === 0 && (
            <EditorEmptyState
              suggestedQueries={suggestedQueries}
              onSelectQuery={setInput}
              onOpenSqlEditor={() => setShowSqlEditor(true)}
            />
          )}
        </AnimatePresence>

        <div>
          <AnimatePresence>
            {messages.map((message, index) => (
              <ChatMessage
                key={message.id}
                message={message}
                index={index}
                copied={copied}
                onCopy={copyToClipboard}
              />
            ))}
          </AnimatePresence>

          {(isSubmitting || audioProcessing) && (
            <TypingIndicator isAudioProcessing={audioProcessing} />
          )}

          <div ref={messagesEndRef} className="h-4" />
        </div>
      </div>

      {/* Input area */}
      <ChatInput
        input={input}
        onInputChange={setInput}
        onSend={onSend}
        onOpenSqlEditor={() => setShowSqlEditor(true)}
        isSubmitting={isSubmitting}
        audioProcessing={audioProcessing}
        isRecording={isRecording}
        recordedBlob={recordedBlob}
        onStartRecording={startRecording}
        onStopRecording={stopRecording}
        onProcessAudio={onProcessAudio}
        onDiscardAudio={discardRecording}
        onClearChat={clearChat}
        hasMessages={messages.length > 0}
      />

      {/* SQL Editor panel */}
      <SqlEditorSheet
        open={showSqlEditor}
        onOpenChange={setShowSqlEditor}
        sqlQuery={sqlQuery}
        onSqlQueryChange={setSqlQuery}
        onSubmit={onSqlSubmit}
        isSubmitting={isSubmitting}
      />
    </div>
  );
}
