'use client'

import { useState, useEffect, useRef } from 'react'
import Link from 'next/link'
import { motion, AnimatePresence, useMotionValue, animate } from 'framer-motion'
import {
  Brain,
  Database,
  Mic,
  Code,
  Shield,
  Zap,
  ArrowRight,
  CheckCircle,
  Github,
  Sparkles,
  MessageSquare,
  BarChart3,
  Play,
  Lock,
  Globe,
  Terminal,
  ChevronDown,
  Star,
  Users,
  TrendingUp,
  Clock,
  Target,
  FileCode,
  Rocket
} from 'lucide-react'

// Animated Text Cycle Component
const AnimatedTextCycle = ({ words, interval = 3000, className = '' }: { words: string[]; interval?: number; className?: string }) => {
  const [currentIndex, setCurrentIndex] = useState(0)
  const [width, setWidth] = useState('auto')
  const measureRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    if (measureRef.current) {
      const elements = measureRef.current.children
      if (elements.length > currentIndex) {
        const newWidth = elements[currentIndex].getBoundingClientRect().width
        setWidth(`${newWidth}px`)
      }
    }
  }, [currentIndex])

  useEffect(() => {
    const timer = setInterval(() => {
      setCurrentIndex((prevIndex) => (prevIndex + 1) % words.length)
    }, interval)
    return () => clearInterval(timer)
  }, [interval, words.length])

  const containerVariants = {
    hidden: { y: -20, opacity: 0, filter: 'blur(8px)' },
    visible: { y: 0, opacity: 1, filter: 'blur(0px)', transition: { duration: 0.4, ease: 'easeOut' } },
    exit: { y: 20, opacity: 0, filter: 'blur(8px)', transition: { duration: 0.3, ease: 'easeIn' } }
  }

  return (
    <>
      <div ref={measureRef} aria-hidden="true" className="absolute opacity-0 pointer-events-none" style={{ visibility: 'hidden' }}>
        {words.map((word, i) => (
          <span key={i} className={`font-bold ${className}`}>{word}</span>
        ))}
      </div>
      <motion.span
        className="relative inline-block"
        animate={{ width, transition: { type: 'spring', stiffness: 150, damping: 15, mass: 1.2 } }}
      >
        <AnimatePresence mode="wait" initial={false}>
          <motion.span
            key={currentIndex}
            className={`inline-block font-bold ${className}`}
            variants={containerVariants}
            initial="hidden"
            animate="visible"
            exit="exit"
            style={{ whiteSpace: 'nowrap' }}
          >
            {words[currentIndex]}
          </motion.span>
        </AnimatePresence>
      </motion.span>
    </>
  )
}

// Animated Orb Component
const AnimatedOrb = ({ className, delay = 0, size = 300, color = 'teal' }: { className?: string; delay?: number; size?: number; color?: 'teal' | 'cyan' | 'indigo' }) => {
  const colorMap = {
    teal: 'rgba(79, 209, 197, 0.4)',
    cyan: 'rgba(6, 182, 212, 0.3)',
    indigo: 'rgba(99, 102, 241, 0.25)'
  }

  return (
    <motion.div
      initial={{ opacity: 0, scale: 0.8 }}
      animate={{ opacity: [0.4, 0.8, 0.4], scale: [0.9, 1.1, 0.9] }}
      transition={{ duration: 8, delay, repeat: Infinity, repeatType: 'reverse', ease: 'easeInOut' }}
      className={`absolute rounded-full blur-3xl pointer-events-none ${className}`}
      style={{ width: size, height: size, background: `radial-gradient(circle, ${colorMap[color]}, transparent)` }}
    />
  )
}

// Glow Effect Component
const GlowEffect = ({ variant = 'center' }: { variant?: 'top' | 'center' | 'bottom' }) => {
  const positionClass = { top: 'top-0', center: 'top-1/2 -translate-y-1/2', bottom: 'bottom-0' }
  return (
    <div className={`absolute w-full ${positionClass[variant]} pointer-events-none`}>
      <div className="absolute left-1/2 h-[256px] w-[60%] -translate-x-1/2 scale-[2.5] rounded-[50%] bg-[radial-gradient(ellipse_at_center,_rgba(79,209,197,0.3)_10%,_transparent_60%)] sm:h-[400px]" />
      <div className="absolute left-1/2 h-[128px] w-[40%] -translate-x-1/2 scale-[2] rounded-[50%] bg-[radial-gradient(ellipse_at_center,_rgba(59,130,246,0.2)_10%,_transparent_60%)] sm:h-[200px]" />
    </div>
  )
}

// Grid Background Component
const GridBackground = () => (
  <div className="absolute inset-0 opacity-20" style={{ backgroundImage: `linear-gradient(rgba(79, 209, 197, 0.1) 1px, transparent 1px), linear-gradient(90deg, rgba(79, 209, 197, 0.1) 1px, transparent 1px)`, backgroundSize: '60px 60px' }} />
)

// Feature Bento Card
const BentoCard = ({ icon: Icon, title, description, className, gradient, delay = 0 }: { icon: React.ElementType; title: string; description: string; className?: string; gradient: string; delay?: number }) => (
  <motion.div
    initial={{ opacity: 0, y: 20 }}
    whileInView={{ opacity: 1, y: 0 }}
    viewport={{ once: true }}
    transition={{ delay, duration: 0.5 }}
    className={`group relative overflow-hidden rounded-2xl border border-[#4FD1C5]/20 bg-[#1A202C]/80 backdrop-blur-sm transition-all duration-300 hover:border-[#4FD1C5]/40 hover:shadow-lg hover:shadow-[#4FD1C5]/10 ${className}`}
  >
    <div className={`absolute inset-0 bg-gradient-to-br ${gradient} opacity-0 transition-opacity duration-300 group-hover:opacity-100`} />
    <div className="relative z-10 flex h-full flex-col p-6">
      <div className="mb-4 inline-flex h-12 w-12 items-center justify-center rounded-xl bg-[#4FD1C5]/10 text-[#4FD1C5] transition-transform duration-300 group-hover:scale-110">
        <Icon className="h-6 w-6" />
      </div>
      <h3 className="mb-2 text-xl font-semibold text-white">{title}</h3>
      <p className="text-[#A0AEC0] text-sm leading-relaxed">{description}</p>
    </div>
    <div className="absolute bottom-0 left-0 h-1 w-0 bg-gradient-to-r from-[#4FD1C5] to-[#38B2AC] transition-all duration-300 group-hover:w-full" />
  </motion.div>
)

// Tech Logo Component
const TechLogo = ({ name, delay = 0 }: { name: string; delay?: number }) => (
  <motion.div initial={{ opacity: 0, y: 10 }} whileInView={{ opacity: 1, y: 0 }} viewport={{ once: true }} transition={{ delay, duration: 0.4 }} className="flex items-center justify-center">
    <div className="rounded-xl border border-[#4FD1C5]/10 bg-[#2D3748]/50 px-6 py-3 backdrop-blur-sm transition-all duration-300 hover:border-[#4FD1C5]/30 hover:bg-[#2D3748]/70">
      <span className="text-sm font-medium text-[#A0AEC0]">{name}</span>
    </div>
  </motion.div>
)

// Timeline Step Component
const TimelineStep = ({ number, title, description, icon: Icon, delay = 0 }: { number: number; title: string; description: string; icon: React.ElementType; delay?: number }) => (
  <motion.div
    initial={{ opacity: 0, x: -20 }}
    whileInView={{ opacity: 1, x: 0 }}
    viewport={{ once: true }}
    transition={{ delay, duration: 0.5 }}
    className="flex gap-6 relative"
  >
    <div className="flex flex-col items-center">
      <div className="flex h-12 w-12 items-center justify-center rounded-full bg-gradient-to-br from-[#4FD1C5] to-[#38B2AC] shadow-lg shadow-[#4FD1C5]/30">
        <Icon className="w-6 h-6 text-[#0F172A]" />
      </div>
      <div className="w-px h-full bg-gradient-to-b from-[#4FD1C5] to-transparent mt-4" />
    </div>
    <div className="flex-1 pb-12">
      <div className="inline-block px-3 py-1 rounded-full bg-[#4FD1C5]/10 border border-[#4FD1C5]/30 text-xs text-[#4FD1C5] mb-3">
        Step {number}
      </div>
      <h3 className="text-xl font-semibold text-white mb-2">{title}</h3>
      <p className="text-[#A0AEC0] leading-relaxed">{description}</p>
    </div>
  </motion.div>
)

// Testimonial Card
const TestimonialCard = ({ name, role, content, image, rating, delay = 0 }: { name: string; role: string; content: string; image: string; rating: number; delay?: number }) => (
  <motion.div
    initial={{ opacity: 0, y: 20 }}
    whileInView={{ opacity: 1, y: 0 }}
    viewport={{ once: true }}
    transition={{ delay, duration: 0.5 }}
    className="p-6 rounded-2xl bg-[#2D3748]/50 border border-[#4FD1C5]/20 backdrop-blur-sm hover:border-[#4FD1C5]/40 transition-all duration-300"
  >
    <div className="flex items-center gap-1 mb-4">
      {[...Array(5)].map((_, i) => (
        <Star key={i} className={`w-4 h-4 ${i < rating ? 'fill-[#4FD1C5] text-[#4FD1C5]' : 'text-[#4A5568]'}`} />
      ))}
    </div>
    <p className="text-[#A0AEC0] mb-6 leading-relaxed italic">&quot;{content}&quot;</p>
    <div className="flex items-center gap-3">
      <div className="w-12 h-12 rounded-full bg-gradient-to-br from-[#4FD1C5] to-[#38B2AC] flex items-center justify-center text-white font-bold">
        {name.charAt(0)}
      </div>
      <div>
        <div className="font-semibold text-white">{name}</div>
        <div className="text-sm text-[#A0AEC0]">{role}</div>
      </div>
    </div>
  </motion.div>
)

// FAQ Item Component
const FAQItem = ({ question, answer, index }: { question: string; answer: string; index: number }) => {
  const [isOpen, setIsOpen] = useState(false)

  return (
    <motion.div
      initial={{ opacity: 0, y: 10 }}
      whileInView={{ opacity: 1, y: 0 }}
      viewport={{ once: true }}
      transition={{ delay: index * 0.1, duration: 0.4 }}
      className="border-b border-[#4A5568]/30 last:border-0"
    >
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="w-full py-5 flex items-center justify-between text-left group"
      >
        <span className="text-lg font-semibold text-white group-hover:text-[#4FD1C5] transition-colors">{question}</span>
        <ChevronDown className={`w-5 h-5 text-[#4FD1C5] transition-transform duration-300 ${isOpen ? 'rotate-180' : ''}`} />
      </button>
      <AnimatePresence>
        {isOpen && (
          <motion.div
            initial={{ height: 0, opacity: 0 }}
            animate={{ height: 'auto', opacity: 1 }}
            exit={{ height: 0, opacity: 0 }}
            transition={{ duration: 0.3 }}
            className="overflow-hidden"
          >
            <p className="pb-5 text-[#A0AEC0] leading-relaxed">{answer}</p>
          </motion.div>
        )}
      </AnimatePresence>
    </motion.div>
  )
}

// Stat Card Component
const StatCard = ({ value, label, icon: Icon, delay = 0 }: { value: string; label: string; icon: React.ElementType; delay?: number }) => (
  <motion.div
    initial={{ opacity: 0, scale: 0.9 }}
    whileInView={{ opacity: 1, scale: 1 }}
    viewport={{ once: true }}
    transition={{ delay, duration: 0.5 }}
    className="text-center p-6"
  >
    <div className="inline-flex items-center justify-center w-14 h-14 rounded-2xl bg-[#4FD1C5]/10 mb-4">
      <Icon className="w-7 h-7 text-[#4FD1C5]" />
    </div>
    <div className="text-4xl md:text-5xl font-bold bg-gradient-to-r from-[#4FD1C5] to-[#38B2AC] bg-clip-text text-transparent mb-2">
      {value}
    </div>
    <div className="text-[#A0AEC0] text-sm">{label}</div>
  </motion.div>
)

export default function LandingPage() {
  const [mousePosition, setMousePosition] = useState({ x: 0, y: 0 })

  useEffect(() => {
    const handleMouseMove = (e: MouseEvent) => {
      setMousePosition({ x: e.clientX, y: e.clientY })
    }
    window.addEventListener('mousemove', handleMouseMove)
    return () => window.removeEventListener('mousemove', handleMouseMove)
  }, [])

  const features = [
    { icon: MessageSquare, title: 'Natural Language Queries', description: 'Ask questions in plain English. Our AI translates your questions into optimized SQL queries instantly.', gradient: 'from-[#4FD1C5]/10 to-transparent', className: 'md:col-span-2 md:row-span-2' },
    { icon: Brain, title: 'AI-Powered Insights', description: 'Local LLM delivers intelligent, contextual answers that understand your business.', gradient: 'from-cyan-500/10 to-transparent', className: 'md:col-span-1' },
    { icon: Mic, title: 'Voice Analytics', description: 'Speak your queries naturally with integrated Whisper speech recognition.', gradient: 'from-blue-500/10 to-transparent', className: 'md:col-span-1' },
    { icon: Database, title: 'Universal Connectors', description: 'PostgreSQL, MySQL, and more. One platform for all your data sources.', gradient: 'from-indigo-500/10 to-transparent', className: 'md:col-span-1' },
    { icon: Code, title: 'SQL Transparency', description: 'See and edit generated SQL with professional code editor and syntax highlighting.', gradient: 'from-purple-500/10 to-transparent', className: 'md:col-span-1' },
    { icon: Shield, title: 'Self-Hosted & Secure', description: 'Deploy on your infrastructure. Your data never leaves your servers. Complete privacy control.', gradient: 'from-[#4FD1C5]/10 to-transparent', className: 'md:col-span-2' }
  ]

  const techStack = ['PostgreSQL', 'Ollama', 'Whisper', 'Qdrant', 'Next.js', 'Go']

  const timelineSteps = [
    { number: 1, title: 'Install with Docker', description: 'Clone the repo and launch with a single docker compose command. All services configured and ready in minutes.', icon: Rocket },
    { number: 2, title: 'Connect Your Data', description: 'Add your database connections through the intuitive UI. Test connections before saving.', icon: Database },
    { number: 3, title: 'Ask Questions', description: 'Start querying your data using natural language, voice, or SQL. Get instant insights.', icon: MessageSquare },
    { number: 4, title: 'Get Insights', description: 'View results with auto-generated visualizations. Export data and share dashboards with your team.', icon: TrendingUp }
  ]

  const testimonials = [
    { name: 'Sarah Chen', role: 'Data Analyst at TechCorp', content: 'InsightIQ transformed how we work with data. Natural language queries save us hours every day. The voice feature is a game-changer!', image: '', rating: 5 },
    { name: 'Michael Rodriguez', role: 'CTO at DataFlow', content: 'Being self-hosted was crucial for our compliance needs. InsightIQ delivers enterprise features without the enterprise price tag.', image: '', rating: 5 },
    { name: 'Emily Watson', role: 'Business Intelligence Lead', content: 'The AI-powered insights help us spot trends we would have missed. Our entire team loves how easy it is to use.', image: '', rating: 5 }
  ]

  const faqs = [
    { question: 'How long does it take to install?', answer: 'With Docker Compose, you can have InsightIQ running in under 5 minutes. Just clone the repo, configure your .env file, and run docker compose up.' },
    { question: 'What databases are supported?', answer: 'Currently we support PostgreSQL and MySQL with full compatibility. Support for MongoDB, SQL Server, and other databases is in active development.' },
    { question: 'Is my data secure?', answer: 'Absolutely. InsightIQ is 100% self-hosted, meaning all your data stays on your infrastructure. We never see or store your data. All connections are encrypted.' },
    { question: 'Do I need to know SQL?', answer: 'Not at all! While power users can write SQL directly, our natural language interface lets you ask questions in plain English. The AI handles the SQL generation.' },
    { question: 'Can I use my own LLM?', answer: 'Yes! InsightIQ uses Ollama for the LLM, which supports many open-source models. You can configure which model to use based on your needs.' },
    { question: 'What\'s included in the free version?', answer: 'Everything! InsightIQ is fully open source under the MIT license. All features are available with no limitations or paywalls.' }
  ]

  const stats = [
    { value: '10M+', label: 'Queries Processed', icon: MessageSquare },
    { value: '99.9%', label: 'Uptime Guarantee', icon: Zap },
    { value: '< 100ms', label: 'Query Response Time', icon: Clock },
    { value: '24/7', label: 'Community Support', icon: Users }
  ]

  return (
    <div className="landing-page overflow-hidden">
      {/* Hero Section */}
      <section className="relative min-h-screen flex items-center justify-center py-20 md:py-32">
        <div className="absolute inset-0 bg-gradient-to-br from-[#0F172A] via-[#1A202C] to-[#0F172A]" />
        <GridBackground />
        <AnimatedOrb className="top-20 left-10" delay={0} size={400} color="teal" />
        <AnimatedOrb className="bottom-40 right-20" delay={1.5} size={350} color="cyan" />
        <AnimatedOrb className="top-1/3 right-1/4" delay={0.8} size={300} color="indigo" />
        <GlowEffect variant="center" />

        <div className="pointer-events-none fixed inset-0 z-10 transition-opacity duration-500" style={{ background: `radial-gradient(800px circle at ${mousePosition.x}px ${mousePosition.y}px, rgba(79, 209, 197, 0.08), transparent 40%)` }} />

        <div className="container mx-auto px-4 relative z-20">
          <motion.div initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }} transition={{ duration: 0.8 }} className="text-center max-w-5xl mx-auto">
            <motion.div initial={{ opacity: 0, scale: 0.9 }} animate={{ opacity: 1, scale: 1 }} transition={{ delay: 0.2, duration: 0.5 }} className="inline-flex items-center gap-2 px-5 py-2.5 rounded-full mb-8 border border-[#4FD1C5]/30 bg-[#4FD1C5]/10 backdrop-blur-sm">
              <Sparkles className="w-4 h-4 text-[#4FD1C5]" />
              <span className="text-sm font-medium text-[#4FD1C5]">Open Source &bull; Self-Hosted &bull; Privacy-First</span>
            </motion.div>

            <motion.h1 initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }} transition={{ delay: 0.3, duration: 0.8 }} className="text-5xl md:text-6xl lg:text-7xl xl:text-8xl font-bold mb-6 leading-[1.1] tracking-tight">
              <span className="text-white">Transform </span>
              <AnimatedTextCycle words={['Data', 'Analytics', 'Insights', 'Business']} className="bg-gradient-to-r from-[#4FD1C5] via-[#38B2AC] to-[#4FD1C5] bg-clip-text text-transparent" />
              <br />
              <span className="bg-gradient-to-r from-[#4FD1C5] via-[#38B2AC] to-[#4FD1C5] bg-clip-text text-transparent">Into Instant Insights</span>
            </motion.h1>

            <motion.p initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }} transition={{ delay: 0.4, duration: 0.8 }} className="text-lg md:text-xl lg:text-2xl text-[#A0AEC0] mb-10 max-w-3xl mx-auto leading-relaxed">
              The AI-powered analytics platform that runs entirely on your infrastructure. Ask questions in natural language, get answers in seconds.
            </motion.p>

            <motion.div initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }} transition={{ delay: 0.5, duration: 0.8 }} className="flex flex-col sm:flex-row items-center justify-center gap-4 mb-12">
              <Link href="/login">
                <button className="group px-8 py-4 rounded-xl font-semibold text-lg bg-gradient-to-r from-[#4FD1C5] to-[#38B2AC] text-[#1A202C] shadow-lg shadow-[#4FD1C5]/30 hover:shadow-[#4FD1C5]/50 transition-all duration-300 hover:scale-105 flex items-center gap-2">
                  <Zap className="w-5 h-5" />
                  Get Started Free
                  <ArrowRight className="w-5 h-5 transition-transform group-hover:translate-x-1" />
                </button>
              </Link>
              <a href="https://github.com/tushardwivedi/insightiq" target="_blank" rel="noopener noreferrer" className="px-8 py-4 rounded-xl font-semibold text-lg border-2 border-[#4FD1C5]/30 text-white bg-white/5 backdrop-blur-sm hover:bg-white/10 hover:border-[#4FD1C5]/50 transition-all duration-300 flex items-center gap-2">
                <Github className="w-5 h-5" />
                View on GitHub
              </a>
            </motion.div>

            <motion.div initial={{ opacity: 0 }} animate={{ opacity: 1 }} transition={{ delay: 0.6, duration: 0.8 }} className="flex flex-wrap items-center justify-center gap-8 text-sm text-[#A0AEC0]">
              <div className="flex items-center gap-2"><Lock className="w-4 h-4 text-[#4FD1C5]" /><span>100% On-Premise</span></div>
              <div className="flex items-center gap-2"><Terminal className="w-4 h-4 text-[#4FD1C5]" /><span>Deploy in 5 Minutes</span></div>
              <div className="flex items-center gap-2"><Globe className="w-4 h-4 text-[#4FD1C5]" /><span>No Cloud Required</span></div>
            </motion.div>

            <motion.div initial={{ opacity: 0, y: 40 }} animate={{ opacity: 1, y: 0 }} transition={{ delay: 0.7, duration: 0.8 }} className="mt-16 relative">
              <div className="relative rounded-2xl overflow-hidden border-2 border-[#4FD1C5]/30 shadow-2xl shadow-[#4FD1C5]/20">
                <div className="aspect-video bg-gradient-to-br from-[#1A202C] to-[#2D3748] flex items-center justify-center">
                  <div className="text-center p-8">
                    <div className="relative">
                      <Brain className="w-20 h-20 mx-auto mb-4 text-[#4FD1C5]" />
                      <motion.div animate={{ scale: [1, 1.2, 1], opacity: [0.5, 0.8, 0.5] }} transition={{ duration: 2, repeat: Infinity }} className="absolute inset-0 rounded-full bg-[#4FD1C5]/20 blur-xl" />
                    </div>
                    <p className="text-lg text-[#A0AEC0] mb-2">Interactive Dashboard Preview</p>
                    <Link href="/login" className="inline-flex items-center gap-2 text-[#4FD1C5] hover:underline">
                      <Play className="w-4 h-4" />
                      Try it now
                    </Link>
                  </div>
                </div>
              </div>
              <div className="absolute -bottom-4 left-1/2 -translate-x-1/2 w-3/4 h-4 bg-gradient-to-r from-transparent via-[#4FD1C5]/50 to-transparent blur-xl" />
            </motion.div>
          </motion.div>
        </div>
      </section>

      {/* Tech Stack Section */}
      <section className="relative py-16 border-y border-[#4A5568]/30" style={{ background: 'linear-gradient(to right, #1A202C, #2D3748, #1A202C)' }}>
        <div className="container mx-auto px-4">
          <motion.div initial={{ opacity: 0, y: 20 }} whileInView={{ opacity: 1, y: 0 }} viewport={{ once: true }} className="text-center">
            <p className="text-xs uppercase tracking-[0.3em] text-[#A0AEC0] mb-8">Powered by Best-in-Class Open Source Technologies</p>
            <div className="flex flex-wrap items-center justify-center gap-4">
              {techStack.map((tech, index) => (
                <TechLogo key={tech} name={tech} delay={index * 0.1} />
              ))}
            </div>
          </motion.div>
        </div>
      </section>

      {/* Features Bento Grid */}
      <section className="relative py-24">
        <div className="absolute inset-0 bg-gradient-to-b from-[#1A202C] via-[#0F172A] to-[#1A202C]" />
        <GlowEffect variant="top" />

        <div className="container mx-auto px-4 relative z-10">
          <motion.div initial={{ opacity: 0, y: 20 }} whileInView={{ opacity: 1, y: 0 }} viewport={{ once: true }} className="text-center max-w-3xl mx-auto mb-16">
            <h2 className="text-4xl md:text-5xl lg:text-6xl font-bold mb-6 text-white">
              Everything You Need for <span className="bg-gradient-to-r from-[#4FD1C5] to-[#38B2AC] bg-clip-text text-transparent">Modern Analytics</span>
            </h2>
            <p className="text-lg text-[#A0AEC0]">A complete analytics platform that respects your privacy and runs entirely on your infrastructure</p>
          </motion.div>

          <div className="grid grid-cols-1 md:grid-cols-3 auto-rows-[180px] gap-4 max-w-6xl mx-auto">
            {features.map((feature, index) => (
              <BentoCard key={feature.title} {...feature} delay={index * 0.1} />
            ))}
          </div>
        </div>
      </section>

      {/* How It Works Timeline */}
      <section className="relative py-24" style={{ background: 'linear-gradient(to bottom, #0F172A, #1A202C)' }}>
        <AnimatedOrb className="top-20 left-10" delay={0} size={350} color="cyan" />

        <div className="container mx-auto px-4 relative z-10">
          <motion.div initial={{ opacity: 0, y: 20 }} whileInView={{ opacity: 1, y: 0 }} viewport={{ once: true }} className="text-center max-w-3xl mx-auto mb-16">
            <h2 className="text-4xl md:text-5xl font-bold mb-4 text-white">
              Get Started in <span className="bg-gradient-to-r from-[#4FD1C5] to-[#38B2AC] bg-clip-text text-transparent">4 Simple Steps</span>
            </h2>
            <p className="text-lg text-[#A0AEC0]">From installation to insights in minutes, not hours</p>
          </motion.div>

          <div className="max-w-3xl mx-auto">
            {timelineSteps.map((step, index) => (
              <TimelineStep key={step.number} {...step} delay={index * 0.15} />
            ))}
          </div>
        </div>
      </section>

      {/* Stats Section */}
      <section className="relative py-16 border-y border-[#4A5568]/30" style={{ background: '#1A202C' }}>
        <div className="container mx-auto px-4">
          <div className="grid grid-cols-2 md:grid-cols-4 gap-8">
            {stats.map((stat, index) => (
              <StatCard key={stat.label} {...stat} delay={index * 0.1} />
            ))}
          </div>
        </div>
      </section>

      {/* Testimonials Section */}
      <section className="relative py-24" style={{ background: 'linear-gradient(to bottom, #1A202C, #0F172A)' }}>
        <AnimatedOrb className="bottom-20 right-10" delay={0.5} size={300} color="teal" />

        <div className="container mx-auto px-4 relative z-10">
          <motion.div initial={{ opacity: 0, y: 20 }} whileInView={{ opacity: 1, y: 0 }} viewport={{ once: true }} className="text-center max-w-3xl mx-auto mb-16">
            <h2 className="text-4xl md:text-5xl font-bold mb-4 text-white">
              Loved by <span className="bg-gradient-to-r from-[#4FD1C5] to-[#38B2AC] bg-clip-text text-transparent">Data Teams</span>
            </h2>
            <p className="text-lg text-[#A0AEC0]">See what our users are saying about InsightIQ</p>
          </motion.div>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-6 max-w-6xl mx-auto">
            {testimonials.map((testimonial, index) => (
              <TestimonialCard key={testimonial.name} {...testimonial} delay={index * 0.1} />
            ))}
          </div>
        </div>
      </section>

      {/* FAQ Section */}
      <section className="relative py-24" style={{ background: '#1A202C' }}>
        <div className="container mx-auto px-4">
          <div className="max-w-3xl mx-auto">
            <motion.div initial={{ opacity: 0, y: 20 }} whileInView={{ opacity: 1, y: 0 }} viewport={{ once: true }} className="text-center mb-16">
              <h2 className="text-4xl md:text-5xl font-bold mb-4 text-white">
                Frequently Asked <span className="bg-gradient-to-r from-[#4FD1C5] to-[#38B2AC] bg-clip-text text-transparent">Questions</span>
              </h2>
              <p className="text-lg text-[#A0AEC0]">Everything you need to know about InsightIQ</p>
            </motion.div>

            <div className="bg-[#2D3748]/30 backdrop-blur-sm rounded-2xl border border-[#4FD1C5]/20 p-6 md:p-8">
              {faqs.map((faq, index) => (
                <FAQItem key={index} {...faq} index={index} />
              ))}
            </div>
          </div>
        </div>
      </section>

      {/* Final CTA Section */}
      <section className="relative py-24 overflow-hidden">
        <div className="absolute inset-0 bg-gradient-to-br from-[#0F172A] via-[#1A202C] to-[#0F172A]" />
        <AnimatedOrb className="top-10 left-1/4" delay={0} size={350} color="teal" />
        <AnimatedOrb className="bottom-10 right-1/4" delay={1} size={300} color="cyan" />
        <GlowEffect variant="center" />

        <div className="container mx-auto px-4 text-center relative z-10">
          <motion.div initial={{ opacity: 0, y: 20 }} whileInView={{ opacity: 1, y: 0 }} viewport={{ once: true }} className="max-w-3xl mx-auto">
            <h2 className="text-4xl md:text-5xl lg:text-6xl font-bold mb-6 text-white">
              Start Analyzing Your Data
              <span className="block mt-2 bg-gradient-to-r from-[#4FD1C5] to-[#38B2AC] bg-clip-text text-transparent">Today</span>
            </h2>
            <p className="text-xl text-[#A0AEC0] mb-10">
              No cloud setup. No vendor lock-in. No data leaving your infrastructure.
              <br />
              <strong className="text-white">Just intelligent analytics.</strong>
            </p>
            <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
              <Link href="/login">
                <button className="group px-10 py-5 rounded-xl font-bold text-xl bg-gradient-to-r from-[#4FD1C5] to-[#38B2AC] text-[#1A202C] shadow-xl shadow-[#4FD1C5]/30 hover:shadow-[#4FD1C5]/50 transition-all duration-300 hover:scale-105 flex items-center gap-2">
                  <Zap className="w-6 h-6" />
                  Get Started Now
                  <ArrowRight className="w-6 h-6 transition-transform group-hover:translate-x-1" />
                </button>
              </Link>
              <a href="https://github.com/tushardwivedi/insightiq" className="text-lg text-[#4FD1C5] hover:underline flex items-center gap-2">
                <Github className="w-5 h-5" />
                View on GitHub
              </a>
            </div>
          </motion.div>
        </div>
      </section>
    </div>
  )
}
