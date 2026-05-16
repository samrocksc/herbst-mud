/* eslint-disable functional/no-class-inheritance, no-restricted-syntax */
import { Component, type ReactNode } from "react";

type Props = Readonly<{
  children: ReactNode
}>

type State = Readonly<{
  hasError: boolean
  error: Error | null
}>

export class ErrorBoundary extends Component<Props, State> {
  state: State = { hasError: false, error: null };

  static getDerivedStateFromError(error: Error): State {
    return { hasError: true, error };
  }

  render() {
    if (!this.state.hasError) return this.props.children;

    return (
      <div className="flex items-center justify-center h-full p-8">
        <div className="text-center">
          <h2 className="text-xl font-bold text-danger mb-2">Something went wrong</h2>
          <p className="text-sm text-text-muted mb-4">{this.state.error?.message}</p>
          <button
            onClick={() => window.location.reload()}
            className="px-4 py-2 bg-primary text-white rounded hover:bg-primary/80"
          >
            Reload
          </button>
        </div>
      </div>
    );
  }
}