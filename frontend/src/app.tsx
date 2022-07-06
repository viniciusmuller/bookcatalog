import Router from 'preact-router'
import { Main } from './pages/main'

export function App() {
  return (
    <Router>
      <Main path="/" />
      <div path="/bye">Bye</div>
    </Router>
  )
}
