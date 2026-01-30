'use client'

import { HeroSection } from '@/components/landing/HeroSection'
import { TechStackSection } from '@/components/landing/TechStackSection'
import { FeaturesSection } from '@/components/landing/FeaturesSection'
import { TimelineSection } from '@/components/landing/TimelineSection'
import { StatsSection } from '@/components/landing/StatsSection'
import { TestimonialsSection } from '@/components/landing/TestimonialsSection'
import { FAQSection } from '@/components/landing/FAQSection'
import { CTASection } from '@/components/landing/CTASection'

export default function LandingPage() {
  return (
    <div className="landing-page overflow-hidden">
      <HeroSection />
      <TechStackSection />
      <FeaturesSection />
      <TimelineSection />
      <StatsSection />
      <TestimonialsSection />
      <FAQSection />
      <CTASection />
    </div>
  )
}
