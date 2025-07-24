import styles from './Button.module.scss'
import clsx from 'clsx'

interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> { // onClick, disabled, type, classname, ...
    option?: 'default' | 'accent' | 'secondary' | 'invisible'
    size?: 'sm' | 'md' | 'lg'
    round?: 'standard' | 'more' | 'max'
    className?: string
    children: React.ReactNode
}

export default function Button({ option = 'default', size = 'md', round = 'standard', className, children, ...props }: ButtonProps) {

    return (
        <button
            className={clsx(styles.button, styles[option], styles[size], styles[round], className)}
            {...props}
        >
            {children}
        </button>
    )
}