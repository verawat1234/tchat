import { faker } from '@faker-js/faker';

// User fixtures
export const createUserFixture = (overrides = {}) => ({
  id: faker.string.uuid(),
  name: faker.person.fullName(),
  email: faker.internet.email(),
  username: faker.internet.username(),
  avatar: faker.image.avatar(),
  bio: faker.person.bio(),
  role: faker.helpers.arrayElement(['user', 'admin', 'moderator']),
  verified: faker.datatype.boolean(),
  createdAt: faker.date.past().toISOString(),
  updatedAt: faker.date.recent().toISOString(),
  ...overrides,
});

// Post fixtures
export const createPostFixture = (overrides = {}) => ({
  id: faker.string.uuid(),
  title: faker.lorem.sentence(),
  content: faker.lorem.paragraphs(),
  excerpt: faker.lorem.paragraph(),
  authorId: faker.string.uuid(),
  author: createUserFixture(),
  tags: faker.helpers.multiple(() => faker.lorem.word(), { count: { min: 1, max: 5 } }),
  likes: faker.number.int({ min: 0, max: 1000 }),
  comments: faker.number.int({ min: 0, max: 100 }),
  published: faker.datatype.boolean(),
  publishedAt: faker.date.past().toISOString(),
  createdAt: faker.date.past().toISOString(),
  updatedAt: faker.date.recent().toISOString(),
  ...overrides,
});

// Product fixtures
export const createProductFixture = (overrides = {}) => ({
  id: faker.string.uuid(),
  name: faker.commerce.productName(),
  description: faker.commerce.productDescription(),
  price: parseFloat(faker.commerce.price()),
  compareAtPrice: faker.datatype.boolean() ? parseFloat(faker.commerce.price()) : null,
  currency: faker.finance.currencyCode(),
  sku: faker.string.alphanumeric(10).toUpperCase(),
  category: faker.commerce.department(),
  brand: faker.company.name(),
  images: faker.helpers.multiple(() => faker.image.url(), { count: { min: 1, max: 5 } }),
  inStock: faker.datatype.boolean(),
  inventory: faker.number.int({ min: 0, max: 1000 }),
  rating: faker.number.float({ min: 1, max: 5, multipleOf: 0.1 }),
  reviews: faker.number.int({ min: 0, max: 500 }),
  features: faker.helpers.multiple(() => faker.commerce.productAdjective(), { count: { min: 3, max: 8 } }),
  variants: [],
  createdAt: faker.date.past().toISOString(),
  updatedAt: faker.date.recent().toISOString(),
  ...overrides,
});

// Message fixtures
export const createMessageFixture = (overrides = {}) => ({
  id: faker.string.uuid(),
  text: faker.lorem.sentence(),
  senderId: faker.string.uuid(),
  sender: createUserFixture(),
  recipientId: faker.string.uuid(),
  conversationId: faker.string.uuid(),
  type: faker.helpers.arrayElement(['text', 'image', 'file', 'audio', 'video']),
  attachments: [],
  read: faker.datatype.boolean(),
  readAt: faker.datatype.boolean() ? faker.date.recent().toISOString() : null,
  edited: faker.datatype.boolean(),
  editedAt: faker.datatype.boolean() ? faker.date.recent().toISOString() : null,
  reactions: [],
  timestamp: faker.date.recent().toISOString(),
  ...overrides,
});

// Notification fixtures
export const createNotificationFixture = (overrides = {}) => ({
  id: faker.string.uuid(),
  type: faker.helpers.arrayElement(['info', 'success', 'warning', 'error', 'mention', 'follow', 'like']),
  title: faker.lorem.sentence(),
  message: faker.lorem.paragraph(),
  icon: faker.helpers.arrayElement(['bell', 'heart', 'comment', 'share', 'user']),
  actionUrl: faker.internet.url(),
  read: faker.datatype.boolean(),
  readAt: faker.datatype.boolean() ? faker.date.recent().toISOString() : null,
  userId: faker.string.uuid(),
  metadata: {},
  createdAt: faker.date.recent().toISOString(),
  ...overrides,
});

// Event fixtures
export const createEventFixture = (overrides = {}) => ({
  id: faker.string.uuid(),
  title: faker.lorem.sentence(),
  description: faker.lorem.paragraphs(),
  startDate: faker.date.future().toISOString(),
  endDate: faker.date.future().toISOString(),
  location: {
    name: faker.location.city(),
    address: faker.location.streetAddress(),
    coordinates: {
      lat: faker.location.latitude(),
      lng: faker.location.longitude(),
    },
  },
  organizer: createUserFixture(),
  attendees: faker.helpers.multiple(() => createUserFixture(), { count: { min: 5, max: 50 } }),
  maxAttendees: faker.number.int({ min: 10, max: 1000 }),
  price: faker.datatype.boolean() ? parseFloat(faker.commerce.price()) : 0,
  currency: faker.finance.currencyCode(),
  tags: faker.helpers.multiple(() => faker.lorem.word(), { count: { min: 1, max: 5 } }),
  image: faker.image.url(),
  status: faker.helpers.arrayElement(['draft', 'published', 'cancelled', 'completed']),
  createdAt: faker.date.past().toISOString(),
  updatedAt: faker.date.recent().toISOString(),
  ...overrides,
});

// Cart fixtures
export const createCartItemFixture = (overrides = {}) => ({
  id: faker.string.uuid(),
  productId: faker.string.uuid(),
  product: createProductFixture(),
  quantity: faker.number.int({ min: 1, max: 10 }),
  price: parseFloat(faker.commerce.price()),
  total: 0, // Will be calculated
  addedAt: faker.date.recent().toISOString(),
  ...overrides,
});

export const createCartFixture = (overrides = {}) => {
  const items = faker.helpers.multiple(
    () => createCartItemFixture(),
    { count: { min: 1, max: 5 } }
  );

  // Calculate totals
  items.forEach(item => {
    item.total = item.price * item.quantity;
  });

  const subtotal = items.reduce((sum, item) => sum + item.total, 0);
  const tax = subtotal * 0.1; // 10% tax
  const shipping = items.length > 0 ? 10 : 0;
  const total = subtotal + tax + shipping;

  return {
    id: faker.string.uuid(),
    userId: faker.string.uuid(),
    items,
    subtotal,
    tax,
    shipping,
    total,
    currency: faker.finance.currencyCode(),
    createdAt: faker.date.past().toISOString(),
    updatedAt: faker.date.recent().toISOString(),
    ...overrides,
  };
};

// Comment fixtures
export const createCommentFixture = (overrides = {}) => ({
  id: faker.string.uuid(),
  content: faker.lorem.paragraph(),
  authorId: faker.string.uuid(),
  author: createUserFixture(),
  postId: faker.string.uuid(),
  parentId: faker.datatype.boolean() ? faker.string.uuid() : null,
  likes: faker.number.int({ min: 0, max: 100 }),
  replies: [],
  edited: faker.datatype.boolean(),
  editedAt: faker.datatype.boolean() ? faker.date.recent().toISOString() : null,
  createdAt: faker.date.past().toISOString(),
  ...overrides,
});

// Settings fixtures
export const createSettingsFixture = (overrides = {}) => ({
  theme: faker.helpers.arrayElement(['light', 'dark', 'auto']),
  language: faker.helpers.arrayElement(['en', 'es', 'fr', 'de', 'ja', 'zh']),
  notifications: {
    email: faker.datatype.boolean(),
    push: faker.datatype.boolean(),
    sms: faker.datatype.boolean(),
    mentions: faker.datatype.boolean(),
    follows: faker.datatype.boolean(),
    likes: faker.datatype.boolean(),
    comments: faker.datatype.boolean(),
  },
  privacy: {
    profileVisibility: faker.helpers.arrayElement(['public', 'friends', 'private']),
    showEmail: faker.datatype.boolean(),
    showPhone: faker.datatype.boolean(),
    allowMessages: faker.helpers.arrayElement(['everyone', 'friends', 'none']),
    allowTags: faker.datatype.boolean(),
  },
  accessibility: {
    fontSize: faker.helpers.arrayElement(['small', 'medium', 'large', 'xlarge']),
    highContrast: faker.datatype.boolean(),
    reducedMotion: faker.datatype.boolean(),
    screenReader: faker.datatype.boolean(),
  },
  ...overrides,
});

// Batch fixture generators
export const createUserFixtures = (count = 5) =>
  faker.helpers.multiple(() => createUserFixture(), { count });

export const createPostFixtures = (count = 10) =>
  faker.helpers.multiple(() => createPostFixture(), { count });

export const createProductFixtures = (count = 20) =>
  faker.helpers.multiple(() => createProductFixture(), { count });

export const createMessageFixtures = (count = 15) =>
  faker.helpers.multiple(() => createMessageFixture(), { count });

export const createNotificationFixtures = (count = 8) =>
  faker.helpers.multiple(() => createNotificationFixture(), { count });

// Test scenario fixtures
export const createAuthenticatedUserScenario = () => ({
  user: createUserFixture({ verified: true }),
  token: faker.string.alphanumeric(40),
  refreshToken: faker.string.alphanumeric(40),
  expiresAt: faker.date.future().toISOString(),
});

export const createChatConversationScenario = () => {
  const users = createUserFixtures(2);
  const messages = faker.helpers.multiple(
    () => createMessageFixture({
      senderId: faker.helpers.arrayElement(users).id,
      sender: faker.helpers.arrayElement(users),
      conversationId: faker.string.uuid(),
    }),
    { count: 10 }
  );

  return {
    conversationId: faker.string.uuid(),
    participants: users,
    messages: messages.sort((a, b) =>
      new Date(a.timestamp).getTime() - new Date(b.timestamp).getTime()
    ),
  };
};

export const createShoppingScenario = () => ({
  user: createUserFixture({ verified: true }),
  products: createProductFixtures(10),
  cart: createCartFixture(),
  orderHistory: faker.helpers.multiple(
    () => ({
      id: faker.string.uuid(),
      items: faker.helpers.multiple(() => createCartItemFixture(), { count: { min: 1, max: 5 } }),
      status: faker.helpers.arrayElement(['pending', 'processing', 'shipped', 'delivered', 'cancelled']),
      total: parseFloat(faker.commerce.price({ min: 50, max: 500 })),
      createdAt: faker.date.past().toISOString(),
    }),
    { count: 5 }
  ),
});

// Utility functions
export const resetFixtureSeed = (seed = 123) => {
  faker.seed(seed);
};

export const createFixtureWithRelations = (type: string, relations = {}) => {
  const factories = {
    user: createUserFixture,
    post: createPostFixture,
    product: createProductFixture,
    message: createMessageFixture,
    notification: createNotificationFixture,
    event: createEventFixture,
    cart: createCartFixture,
    comment: createCommentFixture,
    settings: createSettingsFixture,
  };

  const factory = factories[type];
  if (!factory) {
    throw new Error(`Unknown fixture type: ${type}`);
  }

  return factory(relations);
};