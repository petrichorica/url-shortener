import './output.css'
import React from 'react'
import { useState } from 'react'

const apiURL = "http://goserver:3000"

function Home() {
  const [shorturl, setShorturl] = useState(() => "short URL")
  const [longurl, setLongurl] = useState("")
  const [isCopied, setIsCopied] = useState(false)

  const shortenURL = async () => {
    const requestURL = `${apiURL}/shorten`
    if (!longurl) {
      return
    }
    const form = new FormData()
    form.append('url', longurl)
    const res = await fetch(requestURL, {
      method: 'POST',
      body: form
    })
    try {
      if (res.status === 201) {
        const data = await res.json()
        setShorturl(`${apiURL}/${data.ShortURL}`)
        setIsCopied(false)
      } else {
        const errorMessage = await res.text()
        console.log(errorMessage)
      }
    } catch (error) {
      console.error(error)
    }
  }

  const copyHandler = () => {
    navigator.clipboard.writeText(shorturl)
    setIsCopied(true)
  }

  return (
    <>
      <h1 className="text-4xl text-center mt-10">URL Shortener</h1>
      <div className="flex justify-center mt-10">
        <p className="text-xl">
          Please enter the URL you want to shorten in the input field below.
        </p>
      </div>
      <div className="flex justify-center mt-10">
        <input
          type="text"
          placeholder="Enter URL"
          className="border border-gray-300 p-2 w-96"
          onChange={(e) => setLongurl(e.target.value)}
        />
        <button
          onClick={shortenURL}
          className="bg-blue-500 text-white p-2 ml-2 w-24">Shorten</button>
      </div>
      <div className="flex justify-center mt-10">
        <p className="border border-gray-300 p-2 w-96">{`${shorturl}`}</p>
        <button
          onClick={copyHandler}
          className="bg-blue-500 text-white p-2 ml-2 w-24">{isCopied ? "Copied!" : "Copy"}</button>
      </div>
    </>
  )
}

export default Home