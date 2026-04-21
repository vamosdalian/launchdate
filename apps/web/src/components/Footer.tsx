import { useState } from 'react';
import type { FormEvent } from 'react';
import { Link } from 'react-router-dom';
import { subscribeEmail } from '../services/subscriptionService';

const EMAIL_PATTERN = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
const SUBSCRIBE_ERROR_MESSAGE = 'Unable to subscribe right now. Please try again in a moment.';

const Footer = () => {
  const [email, setEmail] = useState('');
  const [message, setMessage] = useState('');
  const [error, setError] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();

    const normalizedEmail = email.trim();
    if (!normalizedEmail) {
      setError('Please enter an email address.');
      setMessage('');
      return;
    }

    if (!EMAIL_PATTERN.test(normalizedEmail)) {
      setError('Please enter a valid email address.');
      setMessage('');
      return;
    }

    setIsSubmitting(true);
    setError('');
    setMessage('');

    try {
      await subscribeEmail({ email: normalizedEmail });
      setEmail('');
      setMessage('Subscription saved. Check your inbox for a confirmation email with your current status.');
    } catch {
      setError(SUBSCRIBE_ERROR_MESSAGE);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <footer className="border-t border-gray-800 bg-[#111]">
      <div className="container mx-auto px-4 py-12">
        <div className="grid md:grid-cols-3 gap-8">
          <div>
            <div className="mb-4 inline-flex items-center gap-3 text-xs font-semibold uppercase tracking-[0.22em] text-[var(--color-page-fg)]">
              <img src="/rocket.png" alt="LaunchDate" className="h-5 w-5" />
              LaunchDate
            </div>
            <h3 className="mb-3 text-2xl font-bold">Mission control for launch windows, rockets, and space operators.</h3>
            <p className="max-w-md text-gray-400">Track upcoming missions, review recent results, and browse the organizations shaping modern spaceflight.</p>
            <p className="mt-4 text-sm text-gray-500">&copy; 2025 LaunchDate. All Rights Reserved.</p>
          </div>
          <div>
            <h4 className="mb-3 text-lg font-semibold">Quick Links</h4>
            <ul className="space-y-2 text-gray-400">
              <li><Link to="/launches" className="hover:text-white">Launch Dates</Link></li>
              <li><Link to="/rockets" className="hover:text-white">Rockets</Link></li>
              <li><Link to="/companies" className="hover:text-white">Companies</Link></li>
              <li><Link to="/bases" className="hover:text-white">Launch Bases</Link></li>
            </ul>
          </div>
          <div>
            <h4 className="mb-3 text-lg font-semibold">Stay Connected</h4>
            <p className="mb-4 text-gray-400">Subscribe for launch alerts and notable mission updates.</p>
            <form className="space-y-3" onSubmit={handleSubmit}>
              <div className="flex">
                <input
                  type="email"
                  placeholder="Your email"
                  value={email}
                  onChange={(event) => setEmail(event.target.value)}
                  aria-label="Email address"
                  disabled={isSubmitting}
                  required
                  autoComplete="email"
                  className="w-full rounded-l-xl border border-gray-700 bg-gray-800 px-4 py-3 text-sm text-white placeholder:text-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:cursor-not-allowed disabled:opacity-60"
                />
                <button
                  type="submit"
                  disabled={isSubmitting}
                  className="rounded-r-xl bg-blue-600 px-5 font-semibold text-white hover:bg-blue-700 disabled:cursor-not-allowed disabled:bg-blue-800"
                >
                  {isSubmitting ? 'Submitting...' : 'Subscribe'}
                </button>
              </div>
              {error ? (
                <p className="text-sm text-red-400">{error}</p>
              ) : null}
              {message ? (
                <p className="text-sm text-emerald-400">{message}</p>
              ) : null}
            </form>
          </div>
        </div>
      </div>
    </footer>
  );
};

export default Footer;
