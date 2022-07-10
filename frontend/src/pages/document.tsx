import { useEffect, useState } from "preact/hooks"
import { BookCatalogClient } from "../services/bookcatalog-api"
import { AppDocument } from "../types"

export interface DocumentPageProps {
  path: string
  id: string
}

export const DocumentPage = (props: any) => {
  let [document, setDocument] = useState<AppDocument | null>(null)
  const bookCatalogClient = new BookCatalogClient()

  useEffect(() => {
    bookCatalogClient?.getDocument(props.id)
      .then(documentResponse => {
        setDocument(documentResponse)
      })
      .catch(err => console.error("could not load documents", err))
  }, [])

  const deleteDocument = async () => {
    if (!confirm(`Are you sure you want to delete "${document!.title || document!.filename}"?`)) {
      return
    }

    const success = await bookCatalogClient.removeDocument(document!.id)
    if (success) {
      alert("The document was succesfully deleted.")
    }
  }

  return (
    <div style="margin: 2rem">
      {document ?
        <div style="display: flex">
          <img height={700} src={bookCatalogClient.getCoverUrl(document)} />
          <div style="margin-left: 2rem">
            <h2>Title: {document.title || "Unknown"}</h2>
            <h3>Author: {document.author || "Unknown"}</h3>
            <p>Pages: {document.pages || "Unknown"}</p>
            <a href={bookCatalogClient.getDocumentUrl(document)}>Read</a>
            <div style="margin-top: 10px">
              <button
                type='button'
                onClick={deleteDocument}
              >
                Remove book from catalog
              </button>
              <button type='button' onClick={() => history.go(-1)}>Go back</button>
            </div>
          </div>
        </div>
        : <p>Loading...</p>}
    </div>
  )
}
