import { useState } from "preact/hooks"
import { BookCatalogClient } from "../services/bookcatalog-api"

export interface AddDocumentsProps {
  path: string
}

export const AddDocuments = (_props: AddDocumentsProps) => {
  const bookCatalogClient = new BookCatalogClient()
  let [files, setFiles] = useState<File[]>([])

  const submitFiles = async () => {
    const fileTasks = files
      .map(file => bookCatalogClient.uploadDocument(file))
      .map(p =>
        p.then(res => { console.log(res); return res })
          .catch(console.error))

    const importedFiles = await Promise.all(fileTasks)
    const totalFilesImported = importedFiles
      .reduce((acc, v) => v ? acc + 1 : acc, 0)

    alert(`Succesfully imported ${totalFilesImported} out of ${files.length} files.`)
  }

  return (
    <div style="margin: 2rem">
      <label htmlFor="documents-input">Add documents</label>
      <br />
      <input type="file"
        onInput={e => setFiles([...e.currentTarget.files!])} multiple />
      <br />
      <button type="button" onClick={submitFiles} multiple>Submit</button>
    </div>
  )
}
