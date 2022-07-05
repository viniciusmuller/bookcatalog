import { Logo } from './logo'
import { BookCatalogClient } from './services/bookcatalog-api'
import { useEffect, useState } from 'preact/hooks'
import { AppDocument } from './types'

export function App() {
  const bookCatalogClient = new BookCatalogClient()
  let [documents, setDocuments] = useState<AppDocument[]>([])

  useEffect(() => {
    bookCatalogClient.getDocuments()
      .then(documentsResponse => setDocuments(documentsResponse))
  }, [])

  console.log(documents)
  return (
    <div id="catalog">
      <h1>Book Catalog</h1>
      <div id="documents-container">
        {
          documents.map(d =>
            <a key={d.id} href={`http://localhost:8080${d.libraryUrl}`}>
              <img height="300" src={`http://localhost:8080${d.coverUrl}`} title={d.name} />
            </a>
          )
        }
      </div>
    </div>
  )
}
