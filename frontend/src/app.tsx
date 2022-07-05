import { BookCatalogClient } from './services/bookcatalog-api'
import { useEffect, useState } from 'preact/hooks'
import { AppDocument } from './types'
import coverNotFoundImage from './assets/img/cover-not-found.jpg'

export function App() {
  let [documents, setDocuments] = useState<AppDocument[] | null>(null)
  let [serverUrl, setServerUrl] = useState<string>("")
  let [bookCatalogClient, setBookCatalogClient] = useState<BookCatalogClient | null>(null)

  let savedServerUrl = localStorage.getItem('serverUrl')
  if (savedServerUrl && !bookCatalogClient) {
    setServerUrl(savedServerUrl)
    setBookCatalogClient(new BookCatalogClient(savedServerUrl))
  }

  useEffect(() => {
    bookCatalogClient?.getDocuments()
      .then(documentsResponse => {
        localStorage.setItem('serverUrl', serverUrl)
        setDocuments(documentsResponse)
      })
      .catch(err => console.error("could not load documents", err))

  }, [bookCatalogClient])

  const handleInput = (event: any) => {
    setServerUrl(event.target.value);
  };

  const setServer = (event: any) => {
    console.log(serverUrl)
    setBookCatalogClient(new BookCatalogClient(serverUrl))
  };

  return (
    <>
      <label htmlFor='server-url'>Server URL</label>
      <input name='server-url' onChange={handleInput} type='text' size={30} placeholder={serverUrl || 'http://192.168.1.135:8080'} />
      <button type='button' onClick={setServer}>Submit</button>
      <div id="catalog">
        <h1>Book Catalog</h1>
        <input type="text" placeholder="Search for books" />
        <div id="documents-container">
          {documents ?
            documents.map(d =>
              <a key={d.id} href={bookCatalogClient?.getDocumentUrl(d)}>
                <object
                  height="300"
                  data={bookCatalogClient?.getCoverUrl(d)}
                  title={d.name} type="image/jpg"
                >
                  <img height="300" src={coverNotFoundImage} alt="Cover not found" />
                </object>
              </a>
            ) : <p>Loading documents...</p>
          }
        </div>
      </div>
    </>
  )
}
