import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './theme/indigo_cyber.css'
import App from './App.tsx'

// Dynamically inject Google Font dependencies
const link = document.createElement('link')
link.href = 'https://fonts.googleapis.com/css2?family=Outfit:wght@300;400;500;600;700;800&family=Inter:wght@300;400;500;600;700&display=swap'
link.rel = 'stylesheet'
document.head.appendChild(link)

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <App />
  </StrictMode>,
)
