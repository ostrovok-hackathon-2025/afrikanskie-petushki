import { Bouncy } from "ldrs/react";
import "ldrs/react/Bouncy.css";

interface LoaderProps {
  text: string;
}

export default function Loader({ text }: LoaderProps) {
  return (
    <div className="fixed top-0 left-0 w-full h-full flex items-center gap-4 justify-center">
      <div className="font-gain text-xl">{text}</div>
      <Bouncy size="45" speed="1.75" color="black" />
    </div>
  );
}
