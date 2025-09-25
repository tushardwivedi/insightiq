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
      processAudio(file)
    }
  }

  return (
    <div className="bg-white rounded-xl shadow-lg p-6 border border-gray-200">
      <div className="flex items-center gap-3 mb-4">
        <div className="p-2 bg-purple-100 rounded-lg">
          <Mic className="w-5 h-5 text-purple-600" />
        </div>
        <h2 className="text-xl font-semibold text-gray-800">Voice Query</h2>
      </div>

      <div className="space-y-4">
        <div className="flex gap-4">
          <button
            onClick={isRecording ? stopRecording : startRecording}
            disabled={isProcessing}
            className={`flex-1 flex items-center justify-center gap-2 py-3 px-4 rounded-lg transition-all duration-200 ${
              isRecording
                ? 'bg-red-500 hover:bg-red-600 text-white'
                : 'bg-gradient-to-r from-purple-500 to-purple-600 hover:from-purple-600 hover:to-purple-700 text-white'
            } disabled:opacity-50`}
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
            className="flex items-center justify-center gap-2 px-4 py-3 border border-gray-300 rounded-lg hover:bg-gray-50 disabled:opacity-50"
          >
            <Upload className="w-4 h-4" />
            Upload
          </button>
        </div>

        {isRecording && (
          <div className="flex items-center gap-2 p-3 bg-red-50 border border-red-200 rounded-lg">
            <div className="w-3 h-3 bg-red-500 rounded-full animate-pulse"></div>
            <span className="text-red-700 font-medium">Recording... Click "Stop Recording" when done</span>
          </div>
        )}

        {recordedBlob && (
          <div className="p-3 bg-green-50 border border-green-200 rounded-lg">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2">
                <FileAudio className="w-4 h-4 text-green-600" />
                <span className="text-green-700 font-medium">Audio recorded</span>
              </div>
              <button
                onClick={handleRecordedAudio}
                disabled={isProcessing}
                className="px-3 py-1 bg-green-500 text-white rounded text-sm hover:bg-green-600 disabled:opacity-50"
              >
                Process
              </button>
            </div>
          </div>
        )}

        {isProcessing && (
          <div className="flex items-center gap-2 p-3 bg-blue-50 border border-blue-200 rounded-lg">
            <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-blue-500"></div>
            <span className="text-blue-700 font-medium">Processing audio...</span>
          </div>
        )}

        <input
          ref={fileInputRef}
          type="file"
          accept="audio/*"
          onChange={handleFileUpload}
          className="hidden"
        />

        <div className="mt-4 text-sm text-gray-500">
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