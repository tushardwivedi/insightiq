import { useState, useCallback } from "react";

export function useVoiceRecording() {
  const [isRecording, setIsRecording] = useState(false);
  const [recordedBlob, setRecordedBlob] = useState<Blob | null>(null);
  const [mediaRecorder, setMediaRecorder] = useState<MediaRecorder | null>(null);

  const startRecording = useCallback(async () => {
    try {
      const stream = await navigator.mediaDevices.getUserMedia({ audio: true });
      const recorder = new MediaRecorder(stream);
      const newChunks: Blob[] = [];

      recorder.ondataavailable = (e) => {
        if (e.data.size > 0) newChunks.push(e.data);
      };

      recorder.onstop = () => {
        const blob = new Blob(newChunks, { type: "audio/webm" });
        setRecordedBlob(blob);
        stream.getTracks().forEach((t) => t.stop());
      };

      recorder.start(1000);
      setMediaRecorder(recorder);
      setIsRecording(true);
    } catch (error) {
      console.error("Error starting recording:", error);
    }
  }, []);

  const stopRecording = useCallback(() => {
    if (mediaRecorder && isRecording) {
      mediaRecorder.stop();
      setIsRecording(false);
    }
  }, [mediaRecorder, isRecording]);

  const discardRecording = useCallback(() => {
    setRecordedBlob(null);
  }, []);

  return {
    isRecording,
    recordedBlob,
    startRecording,
    stopRecording,
    discardRecording,
  };
}
