'use client'

import { useState, useRef } from 'react'
import { Mic, MicOff, Upload, FileAudio } from 'lucide-react'
import { apiClient } from '@/lib/api'
import { VoiceResponse } from '@/types'

interface Props {
  onResult: (result: VoiceResponse) => void
  onLoading: (loading: boolean) => void
}

export default function VoiceQuerySection({ onResult, onLoading }: Props) {
  const [isRecording, setIsRecording] = useState(false)
  const [isProcessing, setIsProcessing] = useState(false)
  const [recordedBlob, setRecordedBlob] = useState<Blob | null>(null)
  const mediaRecorderRef = useRef<MediaRecorder | null>(null)
  const chunksRef = useRef<Blob[]>([])
  const fileInputRef = useRef<HTMLInputElement>(null)

  const startRecording = async () => {
    try {
      const stream = await navigator.mediaDevices.getUserMedia({ audio: true })
      mediaRecorderRef.current = new MediaRecorder(stream)
      chunksRef.current = []

      mediaRecorderRef.current.ondataavailable = (event) => {
        if (event.data.size > 0) {
          chunksRef.current.push(event.data)
        }
      }

      mediaRecorderRef.current.onstop = () => {
        // Use the MediaRecorder's actual MIME type
        const mimeType = mediaRecorderRef.current?.mimeType || 'audio/webm;codecs=opus'
        const blob = new Blob(chunksRef.current, { type: mimeType })
        setRecordedBlob(blob)
        stream.getTracks().forEach(track => track.stop())
      }

      mediaRecorderRef.current.start(1000)
      setIsRecording(true)
    } catch (error) {
      console.error('Error starting recording:', error)
      alert('Microphone access denied or not available. Please check your browser permissions and try again.')
    }
  }

  const stopRecording = () => {
    if (mediaRecorderRef.current && isRecording) {
      mediaRecorderRef.current.stop()
      setIsRecording(false)
    }
  }

  const processAudio = async (audioFile: File) => {
    setIsProcessing(true)
    onLoading(true)

    try {
      const result = await apiClient.voiceQuery(audioFile)
      onResult(result)
    } catch (error) {
      console.error('Voice query failed:', error)
      // Show error to user by creating an error response
      const errorMessage = error instanceof Error ? error.message : 'Failed to process voice query'
      onResult({
        transcript: 'Error processing audio',
        response: {
          query: 'Voice query failed',
          data: [],
          insights: `Error: ${errorMessage}. Please check if the microphone is working and try again.`,
          timestamp: new Date().toISOString(),
          process_time: '0ms',
          task_id: 'error',
          status: 'failed'
        },
        task_id: 'error',
        process_time: '0ms',
        status: 'failed'
      })
    } finally {
      setIsProcessing(false)
      onLoading(false)
    }
  }

  const handleRecordedAudio = () => {
    if (recordedBlob) {
      // Determine file extension from MIME type
      let extension = 'webm'
      if (recordedBlob.type.includes('wav')) extension = 'wav'
      else if (recordedBlob.type.includes('mp3')) extension = 'mp3'
      else if (recordedBlob.type.includes('m4a')) extension = 'm4a'
      else if (recordedBlob.type.includes('ogg')) extension = 'ogg'

      const file = new File([recordedBlob], `recording.${extension}`, { type: recordedBlob.type })
      processAudio(file)
      setRecordedBlob(null)
    }
  }

  const handleFileUpload = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (file) {
      if (file.size > 10 * 1024 * 1024) {
        alert('File size must be less than 10MB');
        return;
      }

      const allowedTypes = ['audio/wav', 'audio/mp3', 'audio/mpeg', 'audio/m4a', 'audio/webm', 'audio/ogg'];
      if (!allowedTypes.includes(file.type)) {
        alert('Please upload a valid audio file (WAV, MP3, M4A, WebM, OGG)');
        return;
      }

      processAudio(file);
    }
  }

  return (
    <div className="card">
      <div className="flex items-center gap-3 mb-4">
        <div className="p-2 rounded-lg" style={{ background: 'var(--accent-color)', color: 'var(--primary-background)' }}>
          <Mic className="w-5 h-5" />
        </div>
        <h2 className="text-lg font-semibold" style={{ color: 'var(--text-primary)' }}>Voice Query</h2>
      </div>

      <div className="space-y-4">
        <div className="flex gap-4">
          <button
            onClick={isRecording ? stopRecording : startRecording}
            disabled={isProcessing}
            className="flex-1 flex items-center justify-center gap-2 py-3 px-4 rounded-lg transition-all duration-200 disabled:opacity-50"
            style={{
              background: isRecording ? 'var(--error-color)' : 'var(--accent-color)',
              color: 'var(--primary-background)',
              border: 'none'
            }}
          >
            {isRecording ? (
              <>
                <MicOff className="w-4 h-4" />
                Stop Recording
              </>
            ) : (
              <>
                <Mic className="w-4 h-4" />
                Start Recording
              </>
            )}
          </button>

          <button
            onClick={() => fileInputRef.current?.click()}
            disabled={isProcessing}
            className="flex items-center justify-center gap-2 px-4 py-3 rounded-lg disabled:opacity-50 transition-colors"
            style={{
              border: '1px solid var(--border-color)',
              background: 'transparent',
              color: 'var(--text-primary)'
            }}
            onMouseEnter={(e) => e.currentTarget.style.background = 'var(--hover-surface)'}
            onMouseLeave={(e) => e.currentTarget.style.background = 'transparent'}
          >
            <Upload className="w-4 h-4" />
            Upload
          </button>
        </div>

        {isRecording && (
          <div className="flex items-center gap-2 p-3 rounded-lg" style={{ background: 'rgba(252, 129, 129, 0.1)', border: '1px solid var(--error-color)' }}>
            <div className="w-3 h-3 rounded-full animate-pulse" style={{ background: 'var(--error-color)' }}></div>
            <span className="font-medium" style={{ color: 'var(--error-color)' }}>Recording... Click "Stop Recording" when done</span>
          </div>
        )}

        {recordedBlob && (
          <div className="p-3 rounded-lg" style={{ background: 'rgba(104, 211, 145, 0.1)', border: '1px solid var(--success-color)' }}>
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2">
                <FileAudio className="w-4 h-4" style={{ color: 'var(--success-color)' }} />
                <span className="font-medium" style={{ color: 'var(--success-color)' }}>Audio recorded</span>
              </div>
              <button
                onClick={handleRecordedAudio}
                disabled={isProcessing}
                className="px-3 py-1 rounded text-sm transition-colors disabled:opacity-50"
                style={{ background: 'var(--success-color)', color: 'var(--primary-background)' }}
              >
                Process
              </button>
            </div>
          </div>
        )}

        {isProcessing && (
          <div className="flex items-center gap-2 p-3 rounded-lg" style={{ background: 'rgba(79, 209, 197, 0.1)', border: '1px solid var(--accent-color)' }}>
            <div className="animate-spin rounded-full h-4 w-4 border-b-2" style={{ borderColor: 'var(--accent-color)' }}></div>
            <span className="font-medium" style={{ color: 'var(--accent-color)' }}>Processing audio...</span>
          </div>
        )}

        <input
          ref={fileInputRef}
          type="file"
          accept="audio/wav,audio/mp3,audio/mpeg,audio/m4a,audio/webm,audio/ogg"
          onChange={handleFileUpload}
          className="hidden"
        />

        <div className="mt-4 text-sm" style={{ color: 'var(--text-secondary)' }}>
          <p className="font-medium mb-1">Voice Query Instructions:</p>
          <ul className="space-y-1">
            <li>• Click "Start Recording" and speak your query clearly</li>
            <li>• Try: "Show me sales trends for last quarter"</li>
            <li>• Or upload an audio file with your question</li>
            <li>• Supported formats: WAV, MP3, M4A</li>
          </ul>
        </div>
      </div>
    </div>
  )
}