// T046 - FormMessage component with dynamic forms
/**
 * FormMessage Component
 * Displays dynamic forms with various field types, validation, and submission handling
 * Supports conditional fields, file uploads, multi-step forms, and data persistence
 */

import React, { useState, useCallback, useMemo, useRef, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { cn } from '../../lib/utils';
import { MessageData, MessageType, InteractionRequest } from '../../types/MessageData';
import { FormContent, FormField, FormFieldType, FormValidation, FormStatus } from '../../types/FormContent';
import { Button } from '../ui/button';
import { Avatar, AvatarFallback, AvatarImage } from '../ui/avatar';
import { Card, CardContent, CardHeader, CardTitle } from '../ui/card';
import { Badge } from '../ui/badge';
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '../ui/tooltip';
import { Input } from '../ui/input';
import { Textarea } from '../ui/textarea';
import { Checkbox } from '../ui/checkbox';
import { RadioGroup, RadioGroupItem } from '../ui/radio-group';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '../ui/select';
import { Label } from '../ui/label';
import { Progress } from '../ui/progress';
import { Separator } from '../ui/separator';
import {
  Send,
  Save,
  AlertCircle,
  CheckCircle,
  Clock,
  FileText,
  Upload,
  Eye,
  EyeOff,
  RotateCcw,
  ArrowLeft,
  ArrowRight,
  Calendar,
  Hash,
  Type,
  AtSign,
  Phone,
  Link,
  Square,
  Circle
} from 'lucide-react';

// Component Props Interface
interface FormMessageProps {
  message: MessageData & { content: FormContent };
  onInteraction?: (interaction: InteractionRequest) => void;
  onFormSubmit?: (formId: string, data: Record<string, any>) => void;
  onFormSave?: (formId: string, data: Record<string, any>) => void;
  className?: string;
  showAvatar?: boolean;
  showTimestamp?: boolean;
  compactMode?: boolean;
  readonly?: boolean;
  autoSave?: boolean;
  performanceMode?: boolean;
}

// Form Response Type
type FormFieldValue = string | number | boolean | string[] | File[];

interface FormData {
  [fieldId: string]: FormFieldValue;
}

// Animation Variants
const formVariants = {
  initial: { opacity: 0, scale: 0.98, y: 20 },
  animate: { opacity: 1, scale: 1, y: 0 },
  exit: { opacity: 0, scale: 0.98, y: -20 }
};

const fieldVariants = {
  initial: { opacity: 0, x: -10 },
  animate: { opacity: 1, x: 0 },
  exit: { opacity: 0, x: 10 }
};

const stepVariants = {
  initial: { opacity: 0, x: 50 },
  animate: { opacity: 1, x: 0 },
  exit: { opacity: 0, x: -50 }
};

export const FormMessage: React.FC<FormMessageProps> = ({
  message,
  onInteraction,
  onFormSubmit,
  onFormSave,
  className,
  showAvatar = true,
  showTimestamp = true,
  compactMode = false,
  readonly = false,
  autoSave = true,
  performanceMode = false
}) => {
  const formRef = useRef<HTMLDivElement>(null);
  const { content } = message;

  // Form state
  const [formData, setFormData] = useState<FormData>({});
  const [validationErrors, setValidationErrors] = useState<Record<string, string>>({});
  const [currentStep, setCurrentStep] = useState(0);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [isSubmitted, setIsSubmitted] = useState(false);
  const [lastSaved, setLastSaved] = useState<Date | null>(null);

  // Get current step fields
  const currentStepFields = useMemo(() => {
    if (!content.isMultiStep) {
      return content.fields;
    }

    const stepsFields = content.steps?.[currentStep]?.fields || [];
    return content.fields.filter(field => stepsFields.includes(field.id));
  }, [content.fields, content.steps, currentStep, content.isMultiStep]);

  // Calculate form progress
  const formProgress = useMemo(() => {
    const totalFields = content.fields.filter(f => !f.conditional || shouldShowField(f)).length;
    const completedFields = content.fields.filter(f => {
      if (f.conditional && !shouldShowField(f)) return false;
      const value = formData[f.id];
      return value !== undefined && value !== '' && value !== null;
    }).length;

    const stepProgress = content.isMultiStep && content.steps
      ? ((currentStep + 1) / content.steps.length) * 100
      : undefined;

    return {
      completed: completedFields,
      total: totalFields,
      percentage: totalFields > 0 ? (completedFields / totalFields) * 100 : 0,
      stepProgress,
      canProceed: currentStepFields.every(field =>
        !field.required ||
        (formData[field.id] !== undefined && formData[field.id] !== '')
      )
    };
  }, [formData, content.fields, currentStepFields, currentStep, content.steps, content.isMultiStep]);

  // Field visibility logic
  function shouldShowField(field: FormField): boolean {
    if (!field.conditional) return true;

    const condition = field.conditional;
    const dependentValue = formData[condition.dependsOn];

    switch (condition.operator) {
      case 'equals': return dependentValue === condition.value;
      case 'not_equals': return dependentValue !== condition.value;
      case 'contains':
        return Array.isArray(dependentValue)
          ? dependentValue.includes(condition.value)
          : String(dependentValue || '').includes(String(condition.value));
      case 'not_empty': return dependentValue !== undefined && dependentValue !== '';
      case 'empty': return dependentValue === undefined || dependentValue === '';
      default: return true;
    }
  }

  // Validate field
  const validateField = useCallback((field: FormField, value: FormFieldValue): string | null => {
    if (field.required && (value === undefined || value === '' || value === null)) {
      return `${field.label} is required`;
    }

    if (!field.validation || !value) return null;

    const validation = field.validation;

    // String validations
    if (typeof value === 'string') {
      if (validation.minLength && value.length < validation.minLength) {
        return `${field.label} must be at least ${validation.minLength} characters`;
      }
      if (validation.maxLength && value.length > validation.maxLength) {
        return `${field.label} must be no more than ${validation.maxLength} characters`;
      }
      if (validation.pattern) {
        const regex = new RegExp(validation.pattern);
        if (!regex.test(value)) {
          return validation.patternMessage || `${field.label} format is invalid`;
        }
      }
    }

    // Number validations
    if (typeof value === 'number') {
      if (validation.min !== undefined && value < validation.min) {
        return `${field.label} must be at least ${validation.min}`;
      }
      if (validation.max !== undefined && value > validation.max) {
        return `${field.label} must be no more than ${validation.max}`;
      }
    }

    // Array validations (multi-select)
    if (Array.isArray(value)) {
      if (validation.minItems && value.length < validation.minItems) {
        return `Please select at least ${validation.minItems} options`;
      }
      if (validation.maxItems && value.length > validation.maxItems) {
        return `Please select no more than ${validation.maxItems} options`;
      }
    }

    return null;
  }, []);

  // Handle field change
  const handleFieldChange = useCallback((fieldId: string, value: FormFieldValue) => {
    setFormData(prev => ({ ...prev, [fieldId]: value }));

    // Clear validation error for this field
    setValidationErrors(prev => {
      const { [fieldId]: _, ...rest } = prev;
      return rest;
    });

    // Auto-save if enabled
    if (autoSave && !readonly) {
      const saveTimer = setTimeout(() => {
        if (onFormSave) {
          onFormSave(content.id, { ...formData, [fieldId]: value });
          setLastSaved(new Date());
        }
      }, 1000);

      return () => clearTimeout(saveTimer);
    }
  }, [formData, autoSave, readonly, onFormSave, content.id]);

  // Validate form
  const validateForm = useCallback(() => {
    const errors: Record<string, string> = {};
    const fieldsToValidate = content.fields.filter(f => shouldShowField(f));

    fieldsToValidate.forEach(field => {
      const error = validateField(field, formData[field.id]);
      if (error) {
        errors[field.id] = error;
      }
    });

    setValidationErrors(errors);
    return Object.keys(errors).length === 0;
  }, [content.fields, formData, validateField]);

  // Handle step navigation
  const handleNextStep = useCallback(() => {
    if (!content.isMultiStep || !content.steps) return;

    // Validate current step
    const currentStepValid = currentStepFields.every(field => {
      const error = validateField(field, formData[field.id]);
      if (error) {
        setValidationErrors(prev => ({ ...prev, [field.id]: error }));
        return false;
      }
      return true;
    });

    if (currentStepValid && currentStep < content.steps.length - 1) {
      setCurrentStep(currentStep + 1);
    }
  }, [content.isMultiStep, content.steps, currentStepFields, currentStep, formData, validateField]);

  const handlePrevStep = useCallback(() => {
    if (currentStep > 0) {
      setCurrentStep(currentStep - 1);
    }
  }, [currentStep]);

  // Handle form submission
  const handleSubmit = useCallback(async () => {
    if (readonly || isSubmitting) return;

    if (!validateForm()) {
      return;
    }

    setIsSubmitting(true);

    try {
      if (onFormSubmit) {
        await onFormSubmit(content.id, formData);
      }

      if (onInteraction) {
        onInteraction({
          messageId: message.id,
          interactionType: 'form_submit',
          data: { formId: content.id, formData },
          userId: 'current-user',
          timestamp: new Date()
        });
      }

      setIsSubmitted(true);
    } catch (error) {
      console.error('Form submission failed:', error);
    } finally {
      setIsSubmitting(false);
    }
  }, [readonly, isSubmitting, validateForm, onFormSubmit, content.id, formData, onInteraction, message.id]);

  // Handle form reset
  const handleReset = useCallback(() => {
    if (readonly) return;

    setFormData({});
    setValidationErrors({});
    setCurrentStep(0);
    setIsSubmitted(false);
  }, [readonly]);

  // Get field icon
  const getFieldIcon = useCallback((type: FormFieldType) => {
    switch (type) {
      case FormFieldType.TEXT: return <Type className="w-4 h-4" />;
      case FormFieldType.TEXTAREA: return <FileText className="w-4 h-4" />;
      case FormFieldType.EMAIL: return <AtSign className="w-4 h-4" />;
      case FormFieldType.PHONE: return <Phone className="w-4 h-4" />;
      case FormFieldType.URL: return <Link className="w-4 h-4" />;
      case FormFieldType.NUMBER: return <Hash className="w-4 h-4" />;
      case FormFieldType.DATE: return <Calendar className="w-4 h-4" />;
      case FormFieldType.CHECKBOX: return <Square className="w-4 h-4" />;
      case FormFieldType.RADIO: return <Circle className="w-4 h-4" />;
      case FormFieldType.SELECT: return <Circle className="w-4 h-4" />;
      case FormFieldType.MULTI_SELECT: return <Square className="w-4 h-4" />;
      case FormFieldType.FILE: return <Upload className="w-4 h-4" />;
      default: return <Type className="w-4 h-4" />;
    }
  }, []);

  // Render field
  const renderField = useCallback((field: FormField, index: number) => {
    if (!shouldShowField(field)) return null;

    const value = formData[field.id];
    const error = validationErrors[field.id];
    const fieldId = `form-${content.id}-field-${field.id}`;

    const FieldWrapper = ({ children }: { children: React.ReactNode }) => (
      <motion.div
        key={field.id}
        variants={performanceMode ? {} : fieldVariants}
        initial={performanceMode ? {} : "initial"}
        animate={performanceMode ? {} : "animate"}
        exit={performanceMode ? {} : "exit"}
        transition={{ delay: index * 0.05 }}
        className="space-y-2"
      >
        <div className="flex items-center gap-2">
          {getFieldIcon(field.type)}
          <Label htmlFor={fieldId} className="text-sm font-medium">
            {field.label}
            {field.required && <span className="text-destructive ml-1">*</span>}
          </Label>
        </div>

        {field.description && (
          <p className="text-xs text-muted-foreground">{field.description}</p>
        )}

        {children}

        {error && (
          <div className="flex items-center gap-1 text-xs text-destructive">
            <AlertCircle className="w-3 h-3" />
            <span>{error}</span>
          </div>
        )}
      </motion.div>
    );

    switch (field.type) {
      case FormFieldType.TEXT:
      case FormFieldType.EMAIL:
      case FormFieldType.PHONE:
      case FormFieldType.URL:
        return (
          <FieldWrapper key={field.id}>
            <Input
              id={fieldId}
              type={field.type === FormFieldType.EMAIL ? 'email' :
                    field.type === FormFieldType.PHONE ? 'tel' :
                    field.type === FormFieldType.URL ? 'url' : 'text'}
              value={(value as string) || ''}
              onChange={(e) => handleFieldChange(field.id, e.target.value)}
              placeholder={field.placeholder}
              disabled={readonly}
              className={cn(error && "border-destructive focus:border-destructive")}
            />
          </FieldWrapper>
        );

      case FormFieldType.TEXTAREA:
        return (
          <FieldWrapper key={field.id}>
            <Textarea
              id={fieldId}
              value={(value as string) || ''}
              onChange={(e) => handleFieldChange(field.id, e.target.value)}
              placeholder={field.placeholder}
              disabled={readonly}
              rows={3}
              className={cn(error && "border-destructive focus:border-destructive")}
            />
          </FieldWrapper>
        );

      case FormFieldType.NUMBER:
        return (
          <FieldWrapper key={field.id}>
            <Input
              id={fieldId}
              type="number"
              value={(value as number) || ''}
              onChange={(e) => handleFieldChange(field.id, parseInt(e.target.value) || 0)}
              placeholder={field.placeholder}
              disabled={readonly}
              min={field.validation?.min}
              max={field.validation?.max}
              className={cn(error && "border-destructive focus:border-destructive")}
            />
          </FieldWrapper>
        );

      case FormFieldType.DATE:
        return (
          <FieldWrapper key={field.id}>
            <Input
              id={fieldId}
              type="date"
              value={(value as string) || ''}
              onChange={(e) => handleFieldChange(field.id, e.target.value)}
              disabled={readonly}
              className={cn(error && "border-destructive focus:border-destructive")}
            />
          </FieldWrapper>
        );

      case FormFieldType.CHECKBOX:
        return (
          <FieldWrapper key={field.id}>
            <div className="flex items-center space-x-2">
              <Checkbox
                id={fieldId}
                checked={(value as boolean) || false}
                onCheckedChange={(checked) => handleFieldChange(field.id, checked)}
                disabled={readonly}
              />
              <Label htmlFor={fieldId} className="text-sm font-normal cursor-pointer">
                {field.placeholder || 'I agree'}
              </Label>
            </div>
          </FieldWrapper>
        );

      case FormFieldType.RADIO:
        return (
          <FieldWrapper key={field.id}>
            <RadioGroup
              value={(value as string) || ''}
              onValueChange={(newValue) => handleFieldChange(field.id, newValue)}
              disabled={readonly}
            >
              {field.options?.map((option, optionIndex) => (
                <div key={optionIndex} className="flex items-center space-x-2">
                  <RadioGroupItem value={option} id={`${fieldId}-${optionIndex}`} />
                  <Label htmlFor={`${fieldId}-${optionIndex}`} className="text-sm font-normal cursor-pointer">
                    {option}
                  </Label>
                </div>
              ))}
            </RadioGroup>
          </FieldWrapper>
        );

      case FormFieldType.SELECT:
        return (
          <FieldWrapper key={field.id}>
            <Select
              value={(value as string) || ''}
              onValueChange={(newValue) => handleFieldChange(field.id, newValue)}
              disabled={readonly}
            >
              <SelectTrigger className={cn(error && "border-destructive")}>
                <SelectValue placeholder={field.placeholder || 'Select an option'} />
              </SelectTrigger>
              <SelectContent>
                {field.options?.map((option, optionIndex) => (
                  <SelectItem key={optionIndex} value={option}>
                    {option}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </FieldWrapper>
        );

      case FormFieldType.MULTI_SELECT:
        const selectedOptions = (value as string[]) || [];
        return (
          <FieldWrapper key={field.id}>
            <div className="space-y-2 max-h-32 overflow-y-auto">
              {field.options?.map((option, optionIndex) => {
                const isSelected = selectedOptions.includes(option);
                return (
                  <div key={optionIndex} className="flex items-center space-x-2">
                    <Checkbox
                      id={`${fieldId}-${optionIndex}`}
                      checked={isSelected}
                      onCheckedChange={(checked) => {
                        const newValue = checked
                          ? [...selectedOptions, option]
                          : selectedOptions.filter(o => o !== option);
                        handleFieldChange(field.id, newValue);
                      }}
                      disabled={readonly}
                    />
                    <Label htmlFor={`${fieldId}-${optionIndex}`} className="text-sm font-normal cursor-pointer">
                      {option}
                    </Label>
                  </div>
                );
              })}
            </div>
          </FieldWrapper>
        );

      case FormFieldType.FILE:
        return (
          <FieldWrapper key={field.id}>
            <div className="border-2 border-dashed border-border rounded-lg p-4 text-center">
              <Upload className="w-8 h-8 mx-auto text-muted-foreground mb-2" />
              <p className="text-sm text-muted-foreground">
                {field.placeholder || 'Click to upload files'}
              </p>
              <input
                id={fieldId}
                type="file"
                multiple={field.validation?.maxItems !== 1}
                accept={field.validation?.allowedTypes?.join(',')}
                onChange={(e) => {
                  const files = Array.from(e.target.files || []);
                  handleFieldChange(field.id, files);
                }}
                disabled={readonly}
                className="hidden"
              />
              <Button
                variant="outline"
                size="sm"
                onClick={() => document.getElementById(fieldId)?.click()}
                disabled={readonly}
                className="mt-2"
              >
                Choose Files
              </Button>
            </div>
          </FieldWrapper>
        );

      default:
        return null;
    }
  }, [formData, validationErrors, content.id, readonly, handleFieldChange, performanceMode, getFieldIcon]);

  // Performance optimization
  const MotionWrapper = performanceMode ? 'div' : motion.div;
  const motionProps = performanceMode ? {} : {
    variants: formVariants,
    initial: "initial",
    animate: "animate",
    exit: "exit",
    transition: { duration: 0.3, ease: "easeOut" }
  };

  return (
    <TooltipProvider>
      <MotionWrapper
        {...motionProps}
        ref={formRef}
        className={cn(
          "form-message relative group",
          "focus-within:ring-2 focus-within:ring-primary/20 focus-within:ring-offset-2",
          "transition-all duration-200",
          className
        )}
        data-testid={`form-message-${message.id}`}
        data-form-id={content.id}
        role="form"
        aria-label={`Form: ${content.title}`}
      >
        <Card className="form-card">
          <CardHeader className="space-y-3">
            {/* Header */}
            <div className="flex items-start justify-between gap-3">
              <div className="flex items-center gap-3 min-w-0 flex-1">
                {showAvatar && (
                  <Avatar className={cn(compactMode ? "w-8 h-8" : "w-10 h-10")}>
                    <AvatarImage src={`/avatars/${message.senderName.toLowerCase()}.png`} />
                    <AvatarFallback>
                      {message.senderName.substring(0, 2).toUpperCase()}
                    </AvatarFallback>
                  </Avatar>
                )}

                <div className="min-w-0 flex-1">
                  <div className="flex items-center gap-2 flex-wrap">
                    <span className="font-semibold text-foreground truncate">
                      {message.senderName}
                    </span>
                    {message.isOwn && (
                      <Badge variant="secondary" className="text-xs">You</Badge>
                    )}
                  </div>
                  {showTimestamp && (
                    <p className="text-xs text-muted-foreground mt-1">
                      Form • {message.timestamp.toLocaleDateString()}
                    </p>
                  )}
                </div>
              </div>

              <div className="flex items-center gap-2">
                <Badge variant={content.status === FormStatus.ACTIVE ? "default" : "secondary"}>
                  {content.status}
                </Badge>

                {lastSaved && autoSave && (
                  <Badge variant="outline" className="text-xs">
                    <Save className="w-3 h-3 mr-1" />
                    Saved {lastSaved.toLocaleTimeString()}
                  </Badge>
                )}
              </div>
            </div>

            {/* Form title and description */}
            <div className="space-y-2">
              <CardTitle className={cn(
                "text-lg font-bold text-foreground leading-tight",
                compactMode && "text-base"
              )}>
                {content.title}
              </CardTitle>

              {content.description && (
                <p className="text-muted-foreground text-sm leading-relaxed">
                  {content.description}
                </p>
              )}
            </div>

            {/* Progress indicators */}
            {!isSubmitted && (
              <div className="space-y-2">
                {/* Overall progress */}
                <div className="flex justify-between text-xs">
                  <span className="text-muted-foreground">
                    Progress: {formProgress.completed}/{formProgress.total} fields
                  </span>
                  <span className="font-medium">
                    {Math.round(formProgress.percentage)}%
                  </span>
                </div>
                <Progress value={formProgress.percentage} className="h-2" />

                {/* Step progress for multi-step forms */}
                {content.isMultiStep && content.steps && formProgress.stepProgress && (
                  <div className="flex justify-between text-xs">
                    <span className="text-muted-foreground">
                      Step {currentStep + 1} of {content.steps.length}: {content.steps[currentStep].title}
                    </span>
                    <span className="font-medium">
                      {Math.round(formProgress.stepProgress)}%
                    </span>
                  </div>
                )}
              </div>
            )}
          </CardHeader>

          <CardContent className="space-y-4">
            {isSubmitted ? (
              // Success state
              <div className="text-center py-8">
                <CheckCircle className="w-12 h-12 mx-auto mb-3 text-green-600" />
                <h3 className="font-semibold text-foreground mb-2">Form submitted successfully!</h3>
                <p className="text-sm text-muted-foreground">
                  {content.successMessage || 'Thank you for your submission.'}
                </p>
              </div>
            ) : (
              <>
                {/* Form fields */}
                <AnimatePresence mode="wait">
                  <motion.div
                    key={currentStep}
                    variants={performanceMode ? {} : stepVariants}
                    initial={performanceMode ? {} : "initial"}
                    animate={performanceMode ? {} : "animate"}
                    exit={performanceMode ? {} : "exit"}
                    className="space-y-4"
                  >
                    {currentStepFields.map((field, index) => renderField(field, index))}
                  </motion.div>
                </AnimatePresence>

                <Separator />

                {/* Form actions */}
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-2">
                    {content.isMultiStep && currentStep > 0 && (
                      <Button
                        variant="outline"
                        onClick={handlePrevStep}
                        size="sm"
                        className="gap-2"
                      >
                        <ArrowLeft className="w-4 h-4" />
                        Previous
                      </Button>
                    )}

                    <Button
                      variant="ghost"
                      onClick={handleReset}
                      size="sm"
                      className="gap-2"
                      disabled={readonly}
                    >
                      <RotateCcw className="w-4 h-4" />
                      Reset
                    </Button>
                  </div>

                  <div className="flex items-center gap-2">
                    {content.isMultiStep && content.steps && currentStep < content.steps.length - 1 ? (
                      <Button
                        onClick={handleNextStep}
                        disabled={!formProgress.canProceed}
                        size="sm"
                        className="gap-2"
                      >
                        Next
                        <ArrowRight className="w-4 h-4" />
                      </Button>
                    ) : (
                      <Button
                        onClick={handleSubmit}
                        disabled={!formProgress.canProceed || isSubmitting || readonly}
                        size="sm"
                        className="gap-2"
                      >
                        <Send className="w-4 h-4" />
                        {isSubmitting ? 'Submitting...' : 'Submit Form'}
                      </Button>
                    )}
                  </div>
                </div>
              </>
            )}
          </CardContent>
        </Card>

        {/* Performance Debug Info */}
        {process.env.NODE_ENV === 'development' && (
          <div className="absolute top-0 right-0 text-xs text-muted-foreground/50 bg-muted/20 px-1 py-0.5 rounded-bl">
            S: {currentStep + 1}/{content.steps?.length || 1} | F: {formProgress.completed}/{formProgress.total} | P: {performanceMode ? 'ON' : 'OFF'}
          </div>
        )}
      </MotionWrapper>
    </TooltipProvider>
  );
};

// Memoized version for performance optimization
export const MemoizedFormMessage = React.memo(FormMessage, (prevProps, nextProps) => {
  return (
    prevProps.message.id === nextProps.message.id &&
    prevProps.message.timestamp.getTime() === nextProps.message.timestamp.getTime() &&
    prevProps.compactMode === nextProps.compactMode &&
    prevProps.showAvatar === nextProps.showAvatar &&
    prevProps.showTimestamp === nextProps.showTimestamp &&
    prevProps.readonly === nextProps.readonly &&
    prevProps.autoSave === nextProps.autoSave &&
    prevProps.performanceMode === nextProps.performanceMode
  );
});

MemoizedFormMessage.displayName = 'MemoizedFormMessage';

export default FormMessage;