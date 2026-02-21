import React from 'react'
import Layout from './layout/Layout'
// In a real app we would use Toaster here
import { Toaster } from 'react-hot-toast'

function App() {
  return (
    <>
      <Layout />
      <Toaster position="top-right" />
    </>
  )
}

export default App
