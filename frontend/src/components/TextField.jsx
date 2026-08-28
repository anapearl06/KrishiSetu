import Icon from './Icon'

export default function TextField({
  id,
  label,
  labelExtra,
  prefix,
  trailing,
  error = false,
  errorMessage,
  ...inputProps
}) {
  return (
    <div className={`ks-field${error ? ' has-error' : ''}`}>
      <div className="ks-label-row">
        <label className="ks-label" htmlFor={id}>
          {label}
        </label>
        {labelExtra && <span className="ks-label-extra">{labelExtra}</span>}
      </div>

      <div className={prefix ? 'ks-input-group' : 'ks-input-wrap'}>
        {prefix && <span className="ks-input-prefix">{prefix}</span>}
        <input
          id={id}
          name={id}
          className={`ks-input${trailing ? ' ks-has-trailing' : ''}`}
          aria-invalid={error || undefined}
          aria-describedby={error && errorMessage ? `${id}-error` : undefined}
          {...inputProps}
        />
        {trailing && <div className="ks-input-trailing">{trailing}</div>}
      </div>

      {error && errorMessage && (
        <p className="ks-message-error" id={`${id}-error`}>
          <Icon name="info" size={16} />
          {errorMessage}
        </p>
      )}
    </div>
  )
}