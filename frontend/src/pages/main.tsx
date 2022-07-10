import { useEffect, useState } from "preact/hooks"
import { BookCatalogClient } from "../services/bookcatalog-api"
import { AppDocument } from "../types"
import coverNotFoundImage from '../assets/img/cover-not-found.jpg'

export const Main = (_: any) => {
  let [documents, setDocuments] = useState<AppDocument[] | null>(null)
  let [query, _setQuery] = useState<string>("")
  let [page, setPage] = useState<number>(1)
  const bookCatalogClient = new BookCatalogClient()
  const pageSize = 15

  useEffect(() => {
    bookCatalogClient?.getDocuments(query, page, pageSize)
      .then(documentsResponse => {
        setDocuments(documentsResponse)
      })
      .catch(err => console.error("could not load documents", err))
  }, [page, query])

  const setQuery = (query: string) => {
    setPage(1)
    _setQuery(query)
  }

  return (
    <>
      <div style="display: flex; flex-direction: column; align-items: center">
        <h1>Book Catalog</h1>

        <div style="display: flex">
          <input onChange={event => setQuery(event.currentTarget.value)} size={35} type="text" placeholder="Search for books" />
          <button type="button">Search</button>
        </div>
        <a style="margin-top: 10px" href="add-documents"><button type="button">Add documents</button></a>

        <div>
          <p style="display: inline-block; margin-right: 10px">Page {page}</p>
          <button onClick={() => page > 1 && setPage(page - 1)}>Previous page</button>
          <button onClick={() => setPage(page + 1)}>Next page</button>
        </div>
      </div>

      <div id="catalog">
        <div id="documents-container">
          {documents == null && <p>Loading documents...</p>}
          {documents != null && documents.length > 0 &&
            documents.map(d =>
              <a key={d.id} href={`/documents/${d.id}`}>
                <object
                  height="300"
                  data={bookCatalogClient.getCoverUrl(d)}
                  title={d.title || d.filename} type="image/png"
                >
                  <img height="300" src={coverNotFoundImage} alt="Cover not found" />
                </object>
              </a>
            )
          }
          {documents != null && documents.length == 0 &&
            <p>No documents to show.</p>
          }
        </div>
      </div>
    </>
  )
}
