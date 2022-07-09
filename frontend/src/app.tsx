import Router, { Route } from 'preact-router'
import { AddDocuments } from './pages/add-documents'
import { DocumentPage, DocumentPageProps } from './pages/document'
import { Main } from './pages/main'

export function App() {
  return (
    <Router>
      <Main path="/" />
      <DocumentPage path="/documents/:id" />
      <AddDocuments path="/add-documents" />
    </Router>
  )
}
