import { useState, useCallback } from "react";
import { AnalyticsResponse, VoiceResponse } from "@/types";
import { apiClient } from "@/lib/api";

interface Message {
  id: string;
  type: "user" | "assistant";
  content: string;
  sql?: string;
  results?: AnalyticsResponse | VoiceResponse;
  timestamp: Date;
}

export function useChatMessages() {
  const [messages, setMessages] = useState<Message[]>([]);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [audioProcessing, setAudioProcessing] = useState(false);
  const [copied, setCopied] = useState<string | null>(null);

  const handleSend = useCallback(async (input: string) => {
    if (!input.trim() || isSubmitting) return;

    const userMessage: Message = {
      id: Date.now().toString(),
      type: "user",
      content: input,
      timestamp: new Date(),
    };

    setMessages((prev) => [...prev, userMessage]);
    setIsSubmitting(true);

    try {
      const result = await apiClient.textQuery({ query: input });

      const assistantMessage: Message = {
        id: (Date.now() + 1).toString(),
        type: "assistant",
        content: (result as AnalyticsResponse).insights || "Analysis complete",
        results: result as AnalyticsResponse,
        timestamp: new Date(),
      };

      setMessages((prev) => [...prev, assistantMessage]);
    } catch (error) {
      const errorMessage: Message = {
        id: (Date.now() + 1).toString(),
        type: "assistant",
        content: `Error: ${error instanceof Error ? error.message : "Failed to process query"}`,
        timestamp: new Date(),
      };
      setMessages((prev) => [...prev, errorMessage]);
    } finally {
      setIsSubmitting(false);
    }
  }, [isSubmitting]);

  const handleSqlSubmit = useCallback(async (sqlQuery: string) => {
    if (!sqlQuery.trim() || isSubmitting) return;

    const userMessage: Message = {
      id: Date.now().toString(),
      type: "user",
      content: "SQL Query executed",
      sql: sqlQuery,
      timestamp: new Date(),
    };

    setMessages((prev) => [...prev, userMessage]);
    setIsSubmitting(true);

    try {
      const result = await apiClient.sqlQuery({ sql: sqlQuery });

      const assistantMessage: Message = {
        id: (Date.now() + 1).toString(),
        type: "assistant",
        content: (result as AnalyticsResponse).insights || "Analysis complete",
        results: result as AnalyticsResponse,
        timestamp: new Date(),
      };

      setMessages((prev) => [...prev, assistantMessage]);
    } catch (error) {
      const errorMessage: Message = {
        id: (Date.now() + 1).toString(),
        type: "assistant",
        content: `Error: ${error instanceof Error ? error.message : "Failed to execute query"}`,
        timestamp: new Date(),
      };
      setMessages((prev) => [...prev, errorMessage]);
    } finally {
      setIsSubmitting(false);
    }
  }, [isSubmitting]);

  const processAudio = useCallback(async (recordedBlob: Blob) => {
    setAudioProcessing(true);
    try {
      const file = new File([recordedBlob], "recording.webm", { type: "audio/webm" });
      const result = await apiClient.voiceQuery(file);
      const transcript = (result as VoiceResponse).transcript;
      const response = (result as VoiceResponse).response;

      const userMessage: Message = {
        id: Date.now().toString(),
        type: "user",
        content: transcript,
        timestamp: new Date(),
      };

      const assistantMessage: Message = {
        id: (Date.now() + 1).toString(),
        type: "assistant",
        content: response.insights || "Analysis complete",
        results: response,
        timestamp: new Date(),
      };

      setMessages((prev) => [...prev, userMessage, assistantMessage]);
    } catch (error) {
      const errorMessage: Message = {
        id: (Date.now() + 1).toString(),
        type: "assistant",
        content: `Error processing audio: ${error instanceof Error ? error.message : "Unknown error"}`,
        timestamp: new Date(),
      };
      setMessages((prev) => [...prev, errorMessage]);
    } finally {
      setAudioProcessing(false);
    }
  }, []);

  const copyToClipboard = useCallback(async (text: string, id: string) => {
    await navigator.clipboard.writeText(text);
    setCopied(id);
    setTimeout(() => setCopied(null), 2000);
  }, []);

  const clearChat = useCallback(() => {
    setMessages([]);
  }, []);

  return {
    messages,
    isSubmitting,
    audioProcessing,
    copied,
    handleSend,
    handleSqlSubmit,
    processAudio,
    copyToClipboard,
    clearChat,
  };
}
