export function QuestsIcon({ className, stroke }: Readonly<{ className?: string; stroke?: string }>) {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      width="20"
      height="20"
      viewBox="0 0 24 24"
      fill="none"
      stroke={stroke ?? "currentColor"}
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
    >
      <path d="M4 4h16v16H4z" />
      <path d="M9 8h6" />
      <path d="M9 12h4" />
      <path d="M16 16l2 2" />
    </svg>
  );
}