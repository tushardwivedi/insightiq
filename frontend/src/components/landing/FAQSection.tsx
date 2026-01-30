'use client'

import { motion } from 'framer-motion'
import { FAQItem } from './FAQItem'

const faqs = [
  { question: 'How long does it take to install?', answer: 'With Docker Compose, you can have InsightIQ running in under 5 minutes. Just clone the repo, configure your .env file, and run docker compose up.' },
  { question: 'What databases are supported?', answer: 'Currently we support PostgreSQL and MySQL with full compatibility. Support for MongoDB, SQL Server, and other databases is in active development.' },
  { question: 'Is my data secure?', answer: 'Absolutely. InsightIQ is 100% self-hosted, meaning all your data stays on your infrastructure. We never see or store your data. All connections are encrypted.' },
  { question: 'Do I need to know SQL?', answer: 'Not at all! While power users can write SQL directly, our natural language interface lets you ask questions in plain English. The AI handles the SQL generation.' },
  { question: 'Can I use my own LLM?', answer: 'Yes! InsightIQ uses Ollama for the LLM, which supports many open-source models. You can configure which model to use based on your needs.' },
  { question: 'What\'s included in the free version?', answer: 'Everything! InsightIQ is fully open source under the MIT license. All features are available with no limitations or paywalls.' }
]

export function FAQSection() {
  return (
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
  )
}
