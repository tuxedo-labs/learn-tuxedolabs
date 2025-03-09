import React from 'react'
interface AuthLayoutProps {
  children: React.ReactNode
}
const AuthLayout = ({children}: AuthLayoutProps) => {
  return (
    <div className='flex justify-center items-center h-screen'>
    {children}
    </div>
  )
}

export default AuthLayout
