export default function RoleTabs({ active, onChange }) {
  const options = [
    { value: 'farmer', label: 'Farmer' },
    { value: 'buyer', label: 'Buyer' },
  ]

  return (
    <div role="tablist" aria-label="Account type" className="ks-tabs">
      {options.map((option) => (
        <button
          key={option.value}
          type="button"
          role="tab"
          aria-selected={active === option.value}
          className={`ks-tab${active === option.value ? ' ks-tab-active' : ''}`}
          onClick={() => onChange(option.value)}
        >
          {option.label}
        </button>
      ))}
    </div>
  )
}