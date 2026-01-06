'use client'

import { useState, useRef, useEffect } from 'react'
import { Send, Mic, MicOff, Upload, FileAudio } from 'lucide-react'
import { apiClient } from '@/lib/api'
import { AnalyticsResponse, VoiceResponse } from '@/types'
import { Button } from '@/components/ui/button'
import { Textarea } from '@/components/ui/textarea'
import { Card, CardContent } from '@/components/ui/card'

interface Props {
  onResult: (result: AnalyticsResponse | VoiceResponse) => void
  onLoading: (loading: boolean) => void
}

export default function CommandBar({ onResult, onLoading }: Props) {
  const [query, setQuery] = useState('')
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [isRecording, setIsRecording] = useState(false)
  const [isProcessing, setIsProcessing] = useState(false)
  const [recordedBlob, setRecordedBlob] = useState<Blob | null>(null)

  const textareaRef = useRef<HTMLTextAreaElement>(null)
  const mediaRecorderRef = useRef<MediaRecorder | null>(null)
  const chunksRef = useRef<Blob[]>([])
  const fileInputRef = useRef<HTMLInputElement>(null)

  // Auto-resize textarea
  useEffect(() => {
    const textarea = textareaRef.current
    if (textarea) {
      textarea.style.height = 'auto'
      textarea.style.height = Math.min(textarea.scrollHeight, 200) + 'px'
    }
  }, [query])

  const handleTextSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!query.trim() || isSubmitting) return

    setIsSubmitting(true)
    onLoading(true)

    try {
      const result = await apiClient.textQuery({ query })
      onResult(result)
      setQuery('') // Clear input after successful submission
    } catch (error) {
      console.error('Text query failed:', error)
      const errorMessage = error instanceof Error ? error.message : 'Failed to process query'
      onResult({
        query,
        data: [],
        insights: `Error: ${errorMessage}. Please check if the database is accessible and try again.`,
        timestamp: new Date().toISOString(),
        process_time: '0ms',
        task_id: 'error',
        status: 'failed'
      })
    } finally {
      setIsSubmitting(false)
      onLoading(false)
    }
  }

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
      setRecordedBlob(null) // Clear recorded audio after processing
    } catch (error) {
      console.error('Voice query failed:', error)
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
      let extension = 'webm'
      if (recordedBlob.type.includes('wav')) extension = 'wav'
      else if (recordedBlob.type.includes('mp3')) extension = 'mp3'
      else if (recordedBlob.type.includes('m4a')) extension = 'm4a'
      else if (recordedBlob.type.includes('ogg')) extension = 'ogg'

      const file = new File([recordedBlob], `recording.${extension}`, { type: recordedBlob.type })
      processAudio(file)
    }
  }

  const handleFileUpload = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (file) {
      if (file.size > 10 * 1024 * 1024) {
        alert('File size must be less than 10MB')
        return
      }

      const allowedTypes = ['audio/wav', 'audio/mp3', 'audio/mpeg', 'audio/m4a', 'audio/webm', 'audio/ogg']
      if (!allowedTypes.includes(file.type)) {
        alert('Please upload a valid audio file (WAV, MP3, M4A, WebM, OGG)')
        return
      }

      processAudio(file)
    }
  }

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      handleTextSubmit(e as any)
    }
  }

  const isDisabled = isSubmitting || isProcessing || isRecording

  return (
    <>
      {/* Status indicators for voice */}
      {isRecording && (
        <div className="flex items-center gap-2 px-4 py-2 bg-destructive/10 border border-destructive/20 rounded-lg mb-2 mx-auto max-w-2xl">
          <div className="w-3 h-3 bg-destructive rounded-full animate-pulse"></div>
          <span className="text-sm">Recording... Click the mic again to stop</span>
        </div>
      )}

      {recordedBlob && (
        <div className="flex items-center gap-2 px-4 py-2 bg-green-500/10 border border-green-500/20 rounded-lg mb-2 mx-auto max-w-2xl">
          <FileAudio className="w-4 h-4 text-green-600" />
          <span className="text-sm">Audio recorded</span>
          <Button onClick={handleRecordedAudio} disabled={isProcessing} size="sm" variant="outline">
            Process
          </Button>
        </div>
      )}

      {isProcessing && (
        <div className="flex items-center gap-2 px-4 py-2 bg-primary/10 border border-primary/20 rounded-lg mb-2 mx-auto max-w-2xl">
          <div className="w-4 h-4 border-2 border-primary border-t-transparent rounded-full animate-spin"></div>
          <span className="text-sm">Processing audio...</span>
        </div>
      )}

      {/* Main Command Bar */}
      <Card className="max-w-2xl mx-auto">
        <CardContent className="p-3">
          <div className="flex items-end gap-2">
            <Textarea
              ref={textareaRef}
              value={query}
              onChange={(e) => setQuery(e.target.value)}
              onKeyDown={handleKeyDown}
              placeholder="Ask a question or describe the analysis you need..."
              rows={1}
              disabled={isDisabled}
              className="min-h-[40px] resize-none border-0 shadow-none focus-visible:ring-0 p-0"
            />

            <Button
              variant={isRecording ? "destructive" : "outline"}
              size="icon"
              onClick={isRecording ? stopRecording : startRecording}
              disabled={isSubmitting || isProcessing}
              title={isRecording ? "Stop recording" : "Start voice recording"}
            >
              {isRecording ? <MicOff className="w-4 h-4" /> : <Mic className="w-4 h-4" />}
            </Button>

            <Button
              variant="outline"
              size="icon"
              onClick={() => fileInputRef.current?.click()}
              disabled={isDisabled}
              title="Upload audio file"
            >
              <Upload className="w-4 h-4" />
            </Button>

            <Button
              onClick={handleTextSubmit}
              disabled={!query.trim() || isDisabled}
              size="icon"
              title="Send query (Enter)"
            >
              {isSubmitting ? (
                <div className="w-4 h-4 border-2 border-current border-t-transparent rounded-full animate-spin"></div>
              ) : (
                <Send className="w-4 h-4" />
              )}
            </Button>

            <input
              ref={fileInputRef}
              type="file"
              accept="audio/wav,audio/mp3,audio/mpeg,audio/m4a,audio/webm,audio/ogg"
              onChange={handleFileUpload}
              className="hidden"
            />
          </div>
        </CardContent>
      </Card>
    </>
  )
}