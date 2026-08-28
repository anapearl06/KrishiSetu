export default function Icon({ name, size, tone, className = '' }) {
  const style = size ? { fontSize: size } : undefined
  const classes = ['material-symbols-outlined']

  if (tone === 'error') classes.push('ks-icon-error')
  if (className) classes.push(className)

  return (
    <span className={classes.join(' ')} style={style} aria-hidden="true">
      {name}
    </span>
  )
}