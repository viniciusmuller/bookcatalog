import { AppDocument } from "../types";

export class BookCatalogClient {
  constructor(private serverUrl: string) { }

  async getDocuments(): Promise<AppDocument[]> {
    const response = await fetch(`${this.serverUrl}/api/documents`)
    return <AppDocument[]>await response.json()
  }

  getDocumentUrl(doc: AppDocument): string {
    return this.serverUrl + doc.libraryUrl
  }

  // TODO: Handle books without cover (both in front and backend)
  getCoverUrl(doc: AppDocument): string {
    return this.serverUrl + doc.coverUrl
  }
}

