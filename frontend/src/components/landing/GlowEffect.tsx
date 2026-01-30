interface GlowEffectProps {
  variant?: 'top' | 'center' | 'bottom'
}

const positionClass = { top: 'top-0', center: 'top-1/2 -translate-y-1/2', bottom: 'bottom-0' }

export function GlowEffect({ variant = 'center' }: GlowEffectProps) {
  return (
    <div className={`absolute w-full ${positionClass[variant]} pointer-events-none`}>
      <div className="absolute left-1/2 h-[256px] w-[60%] -translate-x-1/2 scale-[2.5] rounded-[50%] bg-[radial-gradient(ellipse_at_center,_rgba(79,209,197,0.3)_10%,_transparent_60%)] sm:h-[400px]" />
      <div className="absolute left-1/2 h-[128px] w-[40%] -translate-x-1/2 scale-[2] rounded-[50%] bg-[radial-gradient(ellipse_at_center,_rgba(59,130,246,0.2)_10%,_transparent_60%)] sm:h-[200px]" />
    </div>
  )
}
