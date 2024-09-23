import './output.css'
import Home from './Home.jsx'
import NotFound from './NotFound.jsx'
import React from 'react'

function App() {
  const renderPage = () => {
    switch (window.location.pathname) {
      case '/':
        return <Home />
      default:
        return <NotFound />
    }
  }

  return renderPage()
}

export default App
