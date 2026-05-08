type FormErrorProps = Readonly<{ message: string }>

export function FormError({ message }: FormErrorProps) {
  return (
    <div className="p-2 bg-danger/10 border border-danger rounded text-danger text-xs">
      {message}
    </div>
  )
}