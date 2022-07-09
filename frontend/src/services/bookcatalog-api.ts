import { AppDocument } from "../types";

export class BookCatalogClient {
  private serverUrl: string

  constructor() {
    this.serverUrl = this.getHost()
  }

  async getDocuments(): Promise<AppDocument[]> {
    const response = await fetch(`${this.serverUrl}/api/documents`)
    return <AppDocument[]>await response.json()
  }

  async uploadDocument(file: File): Promise<AppDocument> {
    const formData = new FormData()
    formData.append("document", file)

    const response = await fetch(`${this.serverUrl}/api/documents`, {
      method: 'POST',
      body: formData
    })
    return <AppDocument>await response.json()
  }

  async removeDocument(id: string): Promise<boolean> {
    const response = await fetch(`${this.serverUrl}/api/documents/${id}`, {
      method: 'DELETE'
    })
    return response.status == 200
  }

  async getDocument(id: string): Promise<AppDocument> {
    const response = await fetch(`${this.serverUrl}/api/documents/${id}`)
    return <AppDocument>await response.json()
  }

  getDocumentUrl(doc: AppDocument): string {
    return this.serverUrl + doc.libraryUrl
  }

  // TODO: Handle books without cover (both in front and backend)
  getCoverUrl(doc: AppDocument): string {
    return this.serverUrl + doc.coverUrl
  }

  private getHost(): string {
    return import.meta.env.DEV ? "http://localhost:8080" : ""
  }
}

