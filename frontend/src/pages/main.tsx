import { useEffect, useState } from "preact/hooks"
import { BookCatalogClient } from "../services/bookcatalog-api"
import { AppDocument } from "../types"
import coverNotFoundImage from '../assets/img/cover-not-found.jpg'

export const Main = (_: any) => {
  let [documents, setDocuments] = useState<AppDocument[] | null>(null)
  let [files, setFiles] = useState<File[]>([])
  const bookCatalogClient = new BookCatalogClient()

  useEffect(() => {
    bookCatalogClient?.getDocuments()
      .then(documentsResponse => {
        setDocuments(documentsResponse)
      })
      .catch(err => console.error("could not load documents", err))
  }, [])

  const submitFiles = () => {
    console.log(files)
  }

  return (
    <>
      <h1>Book Catalog</h1>

      <div style="display: flex;">
        <div>
          <input type="text" placeholder="Search for books" />
        </div>

        <div style="display: flex; flex-direction: column">
          <label htmlFor="documents-input">
            Add documents
          </label>
          <input type="file" onInput={e => setFiles([...e.currentTarget.files!])} multiple />
          <button type="button" onClick={submitFiles} multiple>Submit</button>
        </div>
      </div>

      <div id="catalog">
        <div id="documents-container">
          {documents ?
            documents.map(d =>
              <a key={d.id} href={bookCatalogClient.getDocumentUrl(d)}>
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
