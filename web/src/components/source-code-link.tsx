const SOURCE_URL = "/source/lucas-license-source.tar.gz"

export function SourceCodeLink({ className = "" }: { className?: string }) {
  return (
    <a href={SOURCE_URL} className={className} download>
      Source code
    </a>
  )
}
