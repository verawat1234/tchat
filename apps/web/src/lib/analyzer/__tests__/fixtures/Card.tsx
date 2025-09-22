import React from 'react';
import Button from './Button';

interface CardProps {
  title: string;
  description: string;
  onAction?: () => void;
}

export const Card: React.FC<CardProps> = ({ title, description, onAction }) => {
  return (
    <div className="card">
      <h3>{title}</h3>
      <p>{description}</p>
      {onAction && (
        <Button label="Action" onClick={onAction} />
      )}
    </div>
  );
};

export default Card;