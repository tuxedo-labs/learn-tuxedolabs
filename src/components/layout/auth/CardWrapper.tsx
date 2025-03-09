import { Card, CardContent, CardFooter, CardHeader } from "@/components/ui/card";
import AuthHeader from "../header/AuthHeader";
import React from "react";
import { Button } from "@/components/ui/button";

interface CardAuthWrapperProps {
  label: string;
  title: string;
  backButtonHref: string;
  backButtonLable: string;
  children: React.ReactNode;
}

const CardAuthWrapper = ({ children, label, title, backButtonHref, backButtonLable }: CardAuthWrapperProps) => {
  return (
    <Card className="xl:w-1/4 md:w-1/2 shadow-md">
      <CardHeader>
        <AuthHeader label={label} title={title} />
      </CardHeader>
      <CardContent>
        {children}
      </CardContent>
      <CardFooter>
        <Button variant="link" className="font-normal w-full" size="sm" asChild>
          <a href={backButtonHref}>{backButtonLable}</a>
        </Button>
      </CardFooter>
    </Card>
  )
}

export default CardAuthWrapper
