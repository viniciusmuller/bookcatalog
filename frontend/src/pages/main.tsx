import { useEffect, useState } from "preact/hooks"
import { BookCatalogClient } from "../services/bookcatalog-api"
import { AppDocument } from "../types"
import coverNotFoundImage from '../assets/img/cover-not-found.jpg'

export const Main = (_: any) => {
  let [documents, setDocuments] = useState<AppDocument[] | null>(null)
  const bookCatalogClient = new BookCatalogClient()

  useEffect(() => {
    bookCatalogClient?.getDocuments()
      .then(documentsResponse => {
        setDocuments(documentsResponse)
      })
      .catch(err => console.error("could not load documents", err))
  }, [])

  return (
    <>
      <h1>Book Catalog</h1>

      <div style="display: flex;">
        <input type="text" placeholder="Search for books" />
        <a href="add-documents"><button type="button">Add documents</button></a>
      </div>

      <div id="catalog">
        <div id="documents-container">
          {documents ?
            documents.map(d =>
              <a key={d.id} href={`/documents/${d.id}`}>
                <object
                  height="300"
                  data={bookCatalogClient.getCoverUrl(d)}
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
