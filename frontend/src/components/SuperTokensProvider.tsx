'use client'

import { useEffect } from 'react'
import SuperTokens from 'supertokens-web-js'
import ThirdParty from 'supertokens-web-js/recipe/thirdparty'
import EmailPassword from 'supertokens-web-js/recipe/emailpassword'
import Session from 'supertokens-web-js/recipe/session'

export default function SuperTokensProvider({ children }: { children: React.ReactNode }) {
  useEffect(() => {
    if (typeof window !== 'undefined') {
      SuperTokens.init({
        appInfo: {
          appName: 'InsightIQ',
          apiDomain: process.env.NEXT_PUBLIC_API_DOMAIN || 'http://localhost:8080',
          apiBasePath: '/auth',
        },
        recipeList: [ThirdParty.init(), EmailPassword.init(), Session.init()],
      })
    }
  }, [])

  return <>{children}</>
}
