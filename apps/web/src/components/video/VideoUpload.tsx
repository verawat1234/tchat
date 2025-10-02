import React, { useState, useRef, useCallback } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import type { RootState } from '../../store';
import type {
  VideoUploadRequest,
  VideoUploadResponse,
  ContentRating,
} from '../../types/video';
import { uploadVideo, setUploadProgress } from '../../store/slices/videoSlice';

export interface VideoUploadProps {
  onUploadComplete?: (videoId: string) => void;
  onUploadError?: (error: Error) => void;
  className?: string;
}

const CONTENT_RATINGS: { value: ContentRating; label: string; description: string }[] = [
  { value: 'g', label: 'G - General', description: 'All ages' },
  { value: 'pg', label: 'PG - Parental Guidance', description: 'Some material may not be suitable for children' },
  { value: 'pg13', label: 'PG-13 - Parents Cautioned', description: 'Some material may be inappropriate for children under 13' },
  { value: 'r', label: 'R - Restricted', description: '17+ or accompanied by parent' },
  { value: 'nc17', label: 'NC-17 - Adults Only', description: '18+ only' },
  { value: 'unrated', label: 'Unrated', description: 'Not rated' },
];

const VIDEO_CATEGORIES = [
  'Entertainment',
  'Education',
  'Music',
  'Gaming',
  'Sports',
  'Technology',
  'News',
  'Lifestyle',
  'Travel',
  'Food',
  'Business',
  'Other',
];

export const VideoUpload: React.FC<VideoUploadProps> = ({
  onUploadComplete,
  onUploadError,
  className = '',
}) => {
  const dispatch = useDispatch();
  const fileInputRef = useRef<HTMLInputElement>(null);
  const thumbnailInputRef = useRef<HTMLInputElement>(null);

  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [selectedThumbnail, setSelectedThumbnail] = useState<File | null>(null);
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [tags, setTags] = useState<string[]>([]);
  const [tagInput, setTagInput] = useState('');
  const [contentRating, setContentRating] = useState<ContentRating>('unrated');
  const [category, setCategory] = useState('');
  const [isMonetized, setIsMonetized] = useState(false);
  const [price, setPrice] = useState('');
  const [previewUrl, setPreviewUrl] = useState<string | null>(null);
  const [thumbnailPreviewUrl, setThumbnailPreviewUrl] = useState<string | null>(null);
  const [isUploading, setIsUploading] = useState(false);
  const [uploadError, setUploadError] = useState<string | null>(null);

  // Get upload progress from Redux store
  const uploadProgress = useSelector((state: RootState) =>
    selectedFile ? state.video.uploadProgress[selectedFile.name] || 0 : 0
  );

  // Handle file selection
  const handleFileSelect = useCallback((event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (!file) return;

    // Validate file type
    if (!file.type.startsWith('video/')) {
      setUploadError('Please select a valid video file');
      return;
    }

    // Validate file size (5GB max)
    const maxSize = 5 * 1024 * 1024 * 1024; // 5GB in bytes
    if (file.size > maxSize) {
      setUploadError('File size must be less than 5GB');
      return;
    }

    setSelectedFile(file);
    setUploadError(null);

    // Generate preview
    const url = URL.createObjectURL(file);
    setPreviewUrl(url);
  }, []);

  // Handle thumbnail selection
  const handleThumbnailSelect = useCallback((event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (!file) return;

    // Validate file type
    if (!file.type.startsWith('image/')) {
      setUploadError('Please select a valid image file for thumbnail');
      return;
    }

    setSelectedThumbnail(file);
    setUploadError(null);

    // Generate preview
    const url = URL.createObjectURL(file);
    setThumbnailPreviewUrl(url);
  }, []);

  // Handle tag input
  const handleAddTag = useCallback(() => {
    if (tagInput.trim() && !tags.includes(tagInput.trim())) {
      setTags([...tags, tagInput.trim()]);
      setTagInput('');
    }
  }, [tagInput, tags]);

  // Handle tag removal
  const handleRemoveTag = useCallback((tagToRemove: string) => {
    setTags(tags.filter(tag => tag !== tagToRemove));
  }, [tags]);

  // Handle form submission
  const handleSubmit = useCallback(async (event: React.FormEvent) => {
    event.preventDefault();

    if (!selectedFile) {
      setUploadError('Please select a video file');
      return;
    }

    if (!title.trim()) {
      setUploadError('Please enter a video title');
      return;
    }

    setIsUploading(true);
    setUploadError(null);

    try {
      const uploadRequest: VideoUploadRequest = {
        file: selectedFile,
        title: title.trim(),
        description: description.trim(),
        tags,
        content_rating: contentRating,
        thumbnail: selectedThumbnail || undefined,
        category: category || undefined,
        is_monetized: isMonetized,
        price: isMonetized && price ? parseFloat(price) : undefined,
      };

      // Create FormData for upload
      const formData = new FormData();
      formData.append('file', uploadRequest.file);
      formData.append('title', uploadRequest.title);
      formData.append('description', uploadRequest.description);
      formData.append('tags', uploadRequest.tags.join(','));
      formData.append('content_rating', uploadRequest.content_rating);

      if (uploadRequest.thumbnail) {
        formData.append('thumbnail', uploadRequest.thumbnail);
      }

      if (uploadRequest.category) {
        formData.append('category', uploadRequest.category);
      }

      formData.append('is_monetized', String(uploadRequest.is_monetized));

      if (uploadRequest.price) {
        formData.append('price', String(uploadRequest.price));
      }

      // Upload video with progress tracking
      const response = await fetch('/api/v1/videos', {
        method: 'POST',
        body: formData,
        headers: {
          // Let browser set Content-Type with boundary for multipart/form-data
        },
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || 'Upload failed');
      }

      const data: VideoUploadResponse = await response.json();

      // Success - reset form
      setSelectedFile(null);
      setSelectedThumbnail(null);
      setTitle('');
      setDescription('');
      setTags([]);
      setContentRating('unrated');
      setCategory('');
      setIsMonetized(false);
      setPrice('');
      setPreviewUrl(null);
      setThumbnailPreviewUrl(null);

      onUploadComplete?.(data.video_id);

    } catch (error) {
      console.error('Upload error:', error);
      const errorMessage = error instanceof Error ? error.message : 'Upload failed';
      setUploadError(errorMessage);
      onUploadError?.(error as Error);
    } finally {
      setIsUploading(false);
    }
  }, [
    selectedFile,
    selectedThumbnail,
    title,
    description,
    tags,
    contentRating,
    category,
    isMonetized,
    price,
    onUploadComplete,
    onUploadError,
  ]);

  return (
    <div className={`video-upload-container max-w-4xl mx-auto p-6 ${className}`}>
      <h2 className="text-2xl font-bold mb-6">Upload Video</h2>

      <form onSubmit={handleSubmit} className="space-y-6">
        {/* File Upload Section */}
        <div className="border-2 border-dashed border-gray-300 rounded-lg p-8">
          <input
            ref={fileInputRef}
            type="file"
            accept="video/*"
            onChange={handleFileSelect}
            className="hidden"
          />

          {!selectedFile ? (
            <div className="text-center">
              <svg
                className="mx-auto h-12 w-12 text-gray-400"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"
                />
              </svg>
              <p className="mt-2 text-sm text-gray-600">
                Click to select video file or drag and drop
              </p>
              <p className="text-xs text-gray-500 mt-1">
                Max file size: 5GB
              </p>
              <button
                type="button"
                onClick={() => fileInputRef.current?.click()}
                className="mt-4 px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
              >
                Select Video
              </button>
            </div>
          ) : (
            <div className="space-y-4">
              <div className="flex items-center justify-between">
                <div className="flex items-center space-x-4">
                  {previewUrl && (
                    <video
                      src={previewUrl}
                      className="w-32 h-20 object-cover rounded"
                      controls
                    />
                  )}
                  <div>
                    <p className="font-medium">{selectedFile.name}</p>
                    <p className="text-sm text-gray-500">
                      {(selectedFile.size / (1024 * 1024)).toFixed(2)} MB
                    </p>
                  </div>
                </div>
                <button
                  type="button"
                  onClick={() => {
                    setSelectedFile(null);
                    setPreviewUrl(null);
                  }}
                  className="text-red-500 hover:text-red-700"
                >
                  Remove
                </button>
              </div>

              {isUploading && (
                <div className="w-full bg-gray-200 rounded-full h-2">
                  <div
                    className="bg-blue-500 h-2 rounded-full transition-all"
                    style={{ width: `${uploadProgress}%` }}
                  />
                </div>
              )}
            </div>
          )}
        </div>

        {/* Video Details */}
        <div className="space-y-4">
          {/* Title */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Title *
            </label>
            <input
              type="text"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              placeholder="Enter video title"
              className="w-full px-3 py-2 border border-gray-300 rounded focus:ring-2 focus:ring-blue-500"
              required
            />
          </div>

          {/* Description */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Description
            </label>
            <textarea
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              placeholder="Describe your video"
              rows={4}
              className="w-full px-3 py-2 border border-gray-300 rounded focus:ring-2 focus:ring-blue-500"
            />
          </div>

          {/* Tags */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Tags
            </label>
            <div className="flex gap-2">
              <input
                type="text"
                value={tagInput}
                onChange={(e) => setTagInput(e.target.value)}
                onKeyPress={(e) => e.key === 'Enter' && (e.preventDefault(), handleAddTag())}
                placeholder="Add tags"
                className="flex-1 px-3 py-2 border border-gray-300 rounded focus:ring-2 focus:ring-blue-500"
              />
              <button
                type="button"
                onClick={handleAddTag}
                className="px-4 py-2 bg-gray-200 rounded hover:bg-gray-300"
              >
                Add
              </button>
            </div>
            {tags.length > 0 && (
              <div className="flex flex-wrap gap-2 mt-2">
                {tags.map((tag) => (
                  <span
                    key={tag}
                    className="px-3 py-1 bg-blue-100 text-blue-800 rounded-full text-sm flex items-center gap-2"
                  >
                    {tag}
                    <button
                      type="button"
                      onClick={() => handleRemoveTag(tag)}
                      className="text-blue-600 hover:text-blue-800"
                    >
                      ×
                    </button>
                  </span>
                ))}
              </div>
            )}
          </div>

          {/* Category */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Category
            </label>
            <select
              value={category}
              onChange={(e) => setCategory(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded focus:ring-2 focus:ring-blue-500"
            >
              <option value="">Select category</option>
              {VIDEO_CATEGORIES.map((cat) => (
                <option key={cat} value={cat}>
                  {cat}
                </option>
              ))}
            </select>
          </div>

          {/* Content Rating */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Content Rating *
            </label>
            <select
              value={contentRating}
              onChange={(e) => setContentRating(e.target.value as ContentRating)}
              className="w-full px-3 py-2 border border-gray-300 rounded focus:ring-2 focus:ring-blue-500"
              required
            >
              {CONTENT_RATINGS.map((rating) => (
                <option key={rating.value} value={rating.value}>
                  {rating.label} - {rating.description}
                </option>
              ))}
            </select>
          </div>

          {/* Thumbnail Upload */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Thumbnail (Optional)
            </label>
            <input
              ref={thumbnailInputRef}
              type="file"
              accept="image/*"
              onChange={handleThumbnailSelect}
              className="hidden"
            />
            {thumbnailPreviewUrl ? (
              <div className="flex items-center gap-4">
                <img
                  src={thumbnailPreviewUrl}
                  alt="Thumbnail preview"
                  className="w-32 h-20 object-cover rounded"
                />
                <button
                  type="button"
                  onClick={() => {
                    setSelectedThumbnail(null);
                    setThumbnailPreviewUrl(null);
                  }}
                  className="text-red-500 hover:text-red-700"
                >
                  Remove
                </button>
              </div>
            ) : (
              <button
                type="button"
                onClick={() => thumbnailInputRef.current?.click()}
                className="px-4 py-2 bg-gray-200 rounded hover:bg-gray-300"
              >
                Upload Thumbnail
              </button>
            )}
          </div>

          {/* Monetization */}
          <div className="space-y-2">
            <label className="flex items-center gap-2">
              <input
                type="checkbox"
                checked={isMonetized}
                onChange={(e) => setIsMonetized(e.target.checked)}
                className="rounded"
              />
              <span className="text-sm font-medium text-gray-700">
                Enable Monetization
              </span>
            </label>

            {isMonetized && (
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Price (USD)
                </label>
                <input
                  type="number"
                  value={price}
                  onChange={(e) => setPrice(e.target.value)}
                  placeholder="0.00"
                  min="0"
                  step="0.01"
                  className="w-full px-3 py-2 border border-gray-300 rounded focus:ring-2 focus:ring-blue-500"
                />
              </div>
            )}
          </div>
        </div>

        {/* Error Message */}
        {uploadError && (
          <div className="p-4 bg-red-50 border border-red-200 rounded text-red-700">
            {uploadError}
          </div>
        )}

        {/* Submit Button */}
        <div className="flex gap-4">
          <button
            type="submit"
            disabled={!selectedFile || isUploading}
            className="flex-1 px-6 py-3 bg-blue-500 text-white rounded hover:bg-blue-600 disabled:bg-gray-300 disabled:cursor-not-allowed font-medium"
          >
            {isUploading ? `Uploading... ${uploadProgress}%` : 'Upload Video'}
          </button>
          <button
            type="button"
            onClick={() => {
              setSelectedFile(null);
              setSelectedThumbnail(null);
              setTitle('');
              setDescription('');
              setTags([]);
              setContentRating('unrated');
              setCategory('');
              setIsMonetized(false);
              setPrice('');
              setPreviewUrl(null);
              setThumbnailPreviewUrl(null);
              setUploadError(null);
            }}
            className="px-6 py-3 bg-gray-200 rounded hover:bg-gray-300 font-medium"
          >
            Clear
          </button>
        </div>
      </form>
    </div>
  );
};

export default VideoUpload;