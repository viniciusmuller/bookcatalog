import { AppDocument } from "../types";

export class BookCatalogClient {
  constructor() { }

  async getDocuments(): Promise<AppDocument[]> {
    const response = await fetch("http://localhost:8080/api/documents")
    return <AppDocument[]>await response.json()
  }
}

